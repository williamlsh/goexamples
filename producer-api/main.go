package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	kafkaWriter := getKafkaWriter(kafkaURL, topic)
	defer kafkaWriter.Close()

	// Add handle func for producer.
	http.HandleFunc("/", producerHandler(kafkaWriter))

	// Run the web server.
	fmt.Println("start producer-api ... !!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func producerHandler(kafkaWriter *kafka.Writer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
			Value: body,
		}
		if err := kafkaWriter.WriteMessages(r.Context(), msg); err != nil {
			log.Fatal(err)
		}
	})
}
