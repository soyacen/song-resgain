package service

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
)

var Collection *mongo.Collection

func init() {
	initMongo()
}

func initMongo() {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	Collection = client.Database("girlName").Collection("song")
}
