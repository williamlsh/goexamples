package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaURL = "localhost:9092"
	topic    = "my-topic"
)

func main() {
	kafkaWriter := getKafkaWriter(kafkaURL, topic)
	defer kafkaWriter.Close()

	http.HandleFunc("/", producerHandler(kafkaWriter))

	fmt.Println("start producer-api ... !!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	brokers := strings.Split(kafkaURL, ",")
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Async:        true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				fmt.Printf("Kafka could not write message: %v\n", err)
				return
			}
			for _, msg := range messages {
				fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		},
	}
}

func producerHandler(kafkaWriter *kafka.Writer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		msg := kafka.Message{
			Topic: topic,
			Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
			Value: body,
		}

		ctx, cancer := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancer()

		if err := kafkaWriter.WriteMessages(ctx, msg); err != nil {
			io.WriteString(w, fmt.Sprintf("Kafka could not write message: %v", err))
			fmt.Printf("Kafka could not write message: %v", err)
		}
	}
}
