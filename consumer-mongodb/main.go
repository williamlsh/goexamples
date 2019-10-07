package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// get Mongo db Collection using environment variables.
	mongoURL := os.Getenv("mongoURL")
	dbName := os.Getenv("dbName")
	collectionName := os.Getenv("collectionName")

	collection := getMongoCollection(mongoURL, dbName, collectionName)

	// get kafka reader using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	groupID := os.Getenv("groupID")
	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		res, err := collection.InsertOne(context.Background(), msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", res.InsertedID)
	}
}

func getMongoCollection(mongoURL, dbName, collectionName string) *mongo.Collection {
	clientOpts := options.Client().ApplyURI(mongoURL)
	if err := clientOpts.Validate(); err != nil {
		return nil
	}

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil
	}

	return client.Database(dbName).Collection(collectionName)
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
