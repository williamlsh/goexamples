package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options.
	clientOptions := options.Client().ApplyURI("mongodb://mongoAdmin:secret@localhost:27017")
	if err := clientOptions.Validate(); err != nil {
		log.Fatalf("client options: %v", err)
	}

	// Connect to MongoDB.
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("mongo connect: %s", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Check the connection.
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("client ping: %s", err)
	}

	// Get collection handle.
	collection := client.Database("hack").Collection("box")
	_, err = collection.InsertOne(context.TODO(), bson.M{"fruit": "apple"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"fruit": "apple"},
		bson.D{
			{"$currentDate", bson.D{
				// {"lastModified", true},
				{"lastModified", bson.D{
					{"$type", "timestamp"},
				}},
			}},
			{"$set", bson.M{"fruit": "Banana"}},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := collection.Find(context.TODO(), bson.M{}, options.Find().SetAllowDiskUse(true))
	if err != nil {
		log.Fatal(err)
	}
	var res []struct {
		ID           string    `bson:"_id"`
		Fruit        string    `bson:"fruit"`
		LastModified time.Time `bson:"lastModified"`
	}
	if err := cursor.All(context.TODO(), &res); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %v lastModified: %v\n", res, res[0].LastModified.Local())
}
