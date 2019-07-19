package main

import (
	"goexamples/pb"

	"log"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

var orderServiceURI string

func init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	orderServiceURI = viper.GetString("discovery.orderserivce")
}

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Subscribe("Discovery.OrderService", func(m *nats.Msg) {
		orderServiceDiscovery := pb.ServiceDiscovery{OrderServiceUri: orderServiceURI}
		data, err := proto.Marshal(&orderServiceDiscovery)
		if err == nil {
			conn.Publish(m.Reply, data)
		}
	})
}
