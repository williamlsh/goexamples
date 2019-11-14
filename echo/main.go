package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"net/http"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/williamlsh/goexamples"
)

var echoEndpoint = flag.String("echo_endpoint", "localhost:9000", "endpoint of echo service")

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterMyServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
