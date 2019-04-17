// Reference: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial-part-1-connecting-using-bson-and-crud-operations
package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// Set client options.
	log.Println("Setting client options.")
	clientOptions := options.Client().ApplyURI("mongodb://william:fmx34Mn32z39@ds143326.mlab.com:43326/tutorial")
	if err := clientOptions.Validate(); err != nil {
		log.Fatalf("client options: %s", err)
	}

	// Connect to MongoDB.
	log.Println("Done set client options applying uri, attempting to connect to Mongo server.")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("mongo connect: %s", err)
	}

	// Check the connection.
	log.Println("Initialized new mongo client, attempting to ping Mongo server.")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("client ping: %s", err)
	}

	// Get collection handle.
	log.Println("Connected to MongoDB! Getting collection handle.")
	collection := client.Database("tutorial").Collection("trainers")

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	// Insert documents.
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatalf("insert single doc: %s", err)
	}
	log.Printf("Inserted a single document: %v\n", insertResult.InsertedID)

	insertManyResult, err := collection.InsertMany(context.TODO(), []interface{}{misty, brock})
	if err != nil {
		log.Fatalf("insert multiple docs: %s", err)
	}
	log.Printf("Inserted multiple documents: %v\n", insertManyResult.InsertedIDs)

	// Update documents.
	filter := bson.D{{"name", "Ash"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatalf("update single doc: %s", err)
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find documents.
	// create a value into which the result can be decoded.
	var result Trainer
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatalf("find one doc: %s", err)
	}
	log.Printf("Found a single document: %+v\n", result)

	// Pass these options to the Find method.
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents.
	var results []*Trainer

	// Passing bson.D{{}} as the filter matches all documents in the collection.
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatalf("Find multiple docs: %s", err)
	}

	// Finding multiple documents returns a cursor.
	// Iterating through the cursor allows us to decode documents one at a time.
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded.
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatalf("decode element: %s", err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatalf("iterating through cursor: %s", err)
	}

	// Close the cursor once finished.
	err = cur.Close(context.TODO())
	if err != nil {
		log.Fatalf("close cursor: %s", err)
	}
	log.Printf("Found multiple documents: %+v\n", results)

	// Delete documents.
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatalf("delete documents: %s", err)
	}
	log.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	// Delete entire collection.
	// collection.Drop()

	// Close client connection.
	log.Println("Got collection handle, attempting to disconnect.")
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("close client connection: %s", err)
	}
	log.Println("Connection to MongoDB closed.")
}
