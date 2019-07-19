package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"goexamples/pb"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	msg, err := nc.Request("Discovery.OrderService", nil, 1000*time.Millisecond)
	if err != nil {
		fmt.Println(err)
		return
	}
	if msg == nil {
		return
	}

	orderServiceDiscovery := new(pb.ServiceDiscovery)
	err = proto.Unmarshal(msg.Data, orderServiceDiscovery)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := grpc.Dial(orderServiceDiscovery.OrderServiceUri, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	resp, err := createOrders(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.IsSuccess, resp.Error)

	getOrders(client)
}

func createOrders(client pb.OrderServiceClient) (*pb.OrderResponse, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	order := pb.Order{
		OrderId:   u.String(),
		Status:    "Created",
		CreatedOn: time.Now().Unix(),
		OrderItems: []*pb.Order_OrderItem{
			&pb.Order_OrderItem{
				Code:      "knd100",
				Name:      "Kindle Voyage",
				UnitPrice: 220,
				Quantity:  1,
			},
			&pb.Order_OrderItem{
				Code:      "kc101",
				Name:      "Kindle Voyage SmartShell Case",
				UnitPrice: 10,
				Quantity:  2,
			},
		},
	}
	return client.CreateOrder(context.TODO(), &order)
}

func getOrders(client pb.OrderServiceClient) {
	filter := &pb.OrderFilter{SearchText: ""}
	stream, err := client.GetOrders(context.Background(), filter)
	if err != nil {
		return
	}
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(order)
	}
}
