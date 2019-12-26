package main

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	
	"github.com/williamlsh/goexamples/pb"
)

const (
	exampleScheme      = "example"
	exampleServiceName = "resolver.example.grpc.io"
	
	backendAddr = "localhost:50051"
)

func callUnaryEcho(c pb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := pb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "This is example/name_resolving")
	}
}

func main() {
	passthroughConn, err := grpc.Dial(
		fmt.Sprintf("passthrough:///%s", backendAddr),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer passthroughConn.Close()
	
	fmt.Printf("--- echoing to \"passthrough:///%s\"\n", backendAddr)
	makeRPCs(passthroughConn, 10)
	
	fmt.Println()
	
	exampleConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer exampleConn.Close()
	
	fmt.Printf("--- echoing to \"%s:///%s\"\n", exampleScheme, exampleServiceName)
	makeRPCs(exampleConn, 10)
}

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: {backendAddr},
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string {
	return exampleScheme
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
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*exampleResolver) Close() {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
