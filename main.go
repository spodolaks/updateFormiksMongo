package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("formiks_v2").Collection("submissions")

	// Create filter to match documents where data.projectName exists, is not empty and parentId is not null
	filter := bson.M{
		"data.projectName": bson.M{"$exists": true, "$ne": ""},
		"parentId":         bson.M{"$ne": nil},
	}
	fmt.Println(filter)

	// Create update to unset data.projectName and set data.ProjectType to "    "
	update := bson.M{
		"$unset": bson.M{"data.projectName": ""},
		"$set":   bson.M{"data.projectType": "Historical project"},
	}

	// Perform UpdateMany operation
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Matched %v documents and modified %v documents\n", result.MatchedCount, result.ModifiedCount)

}