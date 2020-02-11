package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	consulapi "github.com/hashicorp/consul/api"

	echo "github.com/williamlsh/goexamples/pb"
)

const (
	consulAddr = ""

	exampleSchema      = "example"
	exampleServiceName = "lb.example.grpc.io"

	rawCfg = `{"loadBalancingPolicy": "round_robin"}`
)

// var addrs = []string{"localhost:50051", "localhost:50052"}

func callUnaryEcho(c echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &echo.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := echo.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/name_resolving")
	}
}

func main() {
	pickfirstConn, err := grpc.Dial(fmt.Sprintf("%s:///%s", exampleSchema, exampleServiceName), grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer pickfirstConn.Close()

	fmt.Println("--- calling helloworld.Greeter/SayHello with pick_first ---")
	makeRPCs(pickfirstConn, 10)

	fmt.Println()

	roundrobinConn, err := grpc.Dial(fmt.Sprintf("%s:///%s", exampleSchema, exampleServiceName),
		grpc.WithDefaultServiceConfig(rawCfg),
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	fmt.Println("--- calling helloworld.Greeter/SayHello with round_robin ---")
	makeRPCs(roundrobinConn, 10)
}

type exampleResolverBuilder struct {
	consulCli *consulapi.Client
}

func (rb *exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs, err := rb.resolveServiceAddrsFromConsul()
	if err != nil {
		return nil, err
	}
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string {
	return exampleSchema
}

func (rb *exampleResolverBuilder) resolveServiceAddrsFromConsul() ([]string, error) {
	catalogs, _, err := rb.consulCli.Catalog().Service(exampleServiceName, "", nil)
	if err != nil {
		return nil, fmt.Errorf("could not call Consul catalog api: %w", err)
	}
	var addrs []string
	for _, catalog := range catalogs {
		addrs = append(addrs, fmt.Sprintf("%s:%d", catalog.Address, catalog.ServicePort))
	}
	fmt.Printf("Resolved addresses: %+v\n", addrs)
	return addrs, nil
}

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{
			Addr: s,
		}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*exampleResolver) Close() {}

func init() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	resolver.Register(&exampleResolverBuilder{consulCli: client})
}
