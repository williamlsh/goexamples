package main

import (
	"context"
	"fmt"
	"log"
	"net"
	
	"google.golang.org/grpc"
	
	"github.com/williamlsh/goexamples/pb"
)

const addr = "localhost:50051"

type ecServer struct {
	pb.UnimplementedEchoServer
	addr string
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", addr)
	}
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &ecServer{addr: addr})
	log.Printf("serving on: %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
