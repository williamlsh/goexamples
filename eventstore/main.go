package main

import (
	"fmt"
	"goexamples/pb"
	"log"

	"github.com/golang/protobuf/proto"

	nats "github.com/nats-io/nats.go"
)

const subject = "Order.>"

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.Subscribe(subject, func(msg *nats.Msg) {
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
