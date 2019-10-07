package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/segmentio/kafka-go"
)

func main() {
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	fmt.Println("start producing ... !!")
	for i := 0; i < 20; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", i)),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}
