package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var repoPath string

func setupPlugins(externalPluginsPath string) error {
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func createTempRepo(ctx context.Context) (string, error) {
	var err error
	repoPath, err = ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	if err := fsrepo.Init(repoPath, cfg); err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	return coreapi.NewCoreAPI(node)
}

func spawnDefault(ctx context.Context) (icore.CoreAPI, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		return nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, err
	}

	return createNode(ctx, defaultPath)
}

func spawnEphemeral(ctx context.Context) (icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	repoPath, err := createTempRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	return createNode(ctx, repoPath)
}

func connectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
	var wg sync.WaitGroup

	peersInfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := peersInfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			peersInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peersInfos))
	for _, peerInfo := range peersInfos {
		peerInfo := peerInfo
		go func() {
			defer wg.Done()
			if err := ipfs.Swarm().Connect(ctx, *peerInfo); err != nil {
				fmt.Printf("failed to connect to %s: %s\n", peerInfo.ID, err)
			}
		}()
	}
	wg.Wait()

	return nil
}

func getUnixfsFile(path string) (files.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		return nil, err
	}

	f, err := files.NewReaderPathFile(path, file, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func main() {
	defer os.RemoveAll(repoPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Spawning node on a temporary repo")

	ipfs, err := spawnEphemeral(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("IPFS node is running")

	fmt.Println("\n-- Adding and getting back files & directories --")

	inputBasePath, err := filepath.Abs("example-folder/")
	if err != nil {
		log.Fatal(err)
	}
	inputPathFile := filepath.Join(inputBasePath, "example.txt")
	inputPathDirectory := filepath.Join(inputBasePath, "test-dir")

	someFile, err := getUnixfsNode(inputPathFile)
	if err != nil {
		log.Fatal(err)
	}

	cidFile, err := ipfs.Unixfs().Add(ctx, someFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added file to IPFS with CID %s\n", cidFile.String())

	someDirectory, err := getUnixfsNode(inputPathDirectory)
	if err != nil {
		log.Fatal(err)
	}

	cidDirectory, err := ipfs.Unixfs().Add(ctx, someDirectory)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added directory to IPFS with CID %s\n", cidDirectory.String())

	outputBasePath := "example-folder"
	outputPathFile := filepath.Join(outputBasePath, strings.Split(cidFile.String(), "/")[2])
	outputPathDirectory := filepath.Join(outputBasePath, strings.Split(cidDirectory.String(), "/")[2])

	rootNodeFile, err := ipfs.Unixfs().Get(ctx, cidFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := files.WriteTo(rootNodeFile, outputPathFile); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got file back from IPFS (IPFS path: %s) and wrote it to %s\n", cidFile.String(), outputPathFile)

	rootNodeDirectory, err := ipfs.Unixfs().Get(ctx, cidDirectory)
	if err != nil {
		log.Fatal(err)
	}
	if err := files.WriteTo(rootNodeDirectory, outputPathDirectory); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got directory back from IPFS (IPFS path: %s) and wrote it to %s\n", cidDirectory.String(), outputPathDirectory)

	bootstrapNodes := []string{
		// IPFS Bootstrapper nodes.
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

		// IPFS Cluster Pinning nodes
		"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
		"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	}
	go connectToPeers(ctx, ipfs, bootstrapNodes)

	exampleCIDStr := "QmUaoioqU7bxezBQZkUcgcSyokatMY71sxsALxQmRRrHrj"
	fmt.Printf("Fetching a file from the network with CID %s\n", exampleCIDStr)

	outputPath := filepath.Join(outputBasePath, exampleCIDStr)
	testCID := icorepath.New(exampleCIDStr)

	rootNode, err := ipfs.Unixfs().Get(ctx, testCID)
	if err != nil {
		log.Fatal(err)
	}
	if err := files.WriteTo(rootNode, outputPath); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote the file to %s\n", outputPath)
}
