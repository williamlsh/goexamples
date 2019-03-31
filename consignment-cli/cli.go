package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "goexamples/consignment-service/proto/consignment"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	addr            = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(filename string) (*pb.Consignment, error) {
	consignment := pb.Consignment{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, err
	}
	return &consignment, nil
}

func main() {
	// Setup a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could nto connect: %v\n", err)
	}
	defer conn.Close()
	client := pb.NewShippingClient(conn)

	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("could not parse file: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("could not create consignment: %v\n", err)
	}
	fmt.Printf("created %t\n", r.Created)

	r, err = client.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("could not get consignment: %v\n", err)
	}
	for _, v := range r.Consignments {
		log.Println(v)
	}
}
