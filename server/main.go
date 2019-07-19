package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"goexamples/pb"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

const (
	port      = ":50051"
	aggregate = "Order"
	event     = "OrderCreated"
)

type server struct{}

func (s server) CreateOrder(ctx context.Context, in *pb.Order) (*pb.OrderResponse, error) {
	// store operation of createing order.
	// ...
	go publishOrderCreated(in)
	return &pb.OrderResponse{IsSuccess: true}, nil
}

func (s server) GetOrders(filter *pb.OrderFilter, stream pb.OrderService_GetOrdersServer) error {
	// store operation of getting orders.
	// ...
	orders := []*pb.Order{nil, nil, nil}
	for _, order := range orders {
		if err := stream.Send(order); err != nil {
			return err
		}
	}
	return nil
}

func publishOrderCreated(order *pb.Order) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer nc.Close()

	u, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}

	eventData, err := json.Marshal(order)
	event := pb.EventStore{
		AggregateId:   order.OrderId,
		AggregateType: aggregate,
		EventId:       u.String(),
		EventType:     event,
		EventData:     string(eventData),
	}
	subj := "Order.OrderCreated"
	data, err := proto.Marshal(&event)
	if err != nil {
		fmt.Println(err)
		return
	}
	nc.Publish(subj, data)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, server{})
	s.Serve(lis)
}
