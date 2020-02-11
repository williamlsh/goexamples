package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"

	pb "google.golang.org/grpc/examples/features/proto/echo"
)

const consulAddr = ":8500"

var (
	addrs = []string{":50051"}
)

type ecServer struct {
	pb.UnimplementedEchoServer
	addr string
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &ecServer{addr: addr})
	log.Printf("serving on %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	address, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	registration := &consulapi.AgentServiceRegistration{
		Name:    "lb.example.grpc.io",
		Port:    50051,
		Address: address,
	}
	if err := client.Agent().ServiceRegister(registration); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}
	wg.Wait()
}
