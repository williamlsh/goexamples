package main

import (
	"fmt"
	"log"

	"goexamples/pb"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
)

const (
	queue   = "Order.OrdersCreatedQueue"
	subject = "Order.OrderCreated"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		eventStore := new(pb.EventStore)
		err := proto.Unmarshal(msg.Data, eventStore)
		if err != nil {
			fmt.Println(err)
			return
		}
		// store operation of event.
		// ...
	})
}
