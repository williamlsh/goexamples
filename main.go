package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"goexample/pb"
)

func main() {
	log.Fatal(run())
}

type myService struct{}

func (m *myService) Echo(ctx context.Context, s *pb.StringMessage) (*pb.StringMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

func newServer() *myService {
	return new(myService)
}

func run() error {
	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterEchoServiceServer(s, newServer())
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, "127.0.0.1:8080", opts)
	if err != nil {
		return err
	}

	fmt.Println("Serving HTTP at 127.0.0.1:8081")
	return http.ListenAndServe("127.0.0.1:8081", mux)
}
