package main

import (
	"fmt"
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	// Zookeeper client port.
	hosts := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	hostPro := new(zk.DNSHostProvider)
	err := hostPro.Init(hosts)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = hostPro.Next()

	hostPro.Connected()

	opt := zk.WithEventCallback(func(e zk.Event) {
		fmt.Printf("Event: %+v\n", e)
	})

	conn, _, err := zk.Connect(hosts, 5*time.Second, opt)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var (
		path = "/test"
		data = []byte("hello world")
		flag = zk.FlagEphemeral
		acls = zk.WorldACL(zk.PermAll)
	)

	_, _, _, err = conn.ExistsW(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Create(path, data, int32(flag), acls)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Delete(path, int32(1))
	if err != nil {
		log.Fatal(err)
	}
}
