package main

import (
	"goexamples/pb"
	"net/http"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --twirp_out=. --twirp_opt=paths=source_relative pb/service.proto

func main() {
	server := &Server{}
	twirpHandler := pb.NewHaberdasherServer(server)

	http.ListenAndServe(":8080", twirpHandler)
}
