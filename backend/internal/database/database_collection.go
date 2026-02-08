package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// this function will connect to the databse with refrence to the enviornment of prod or dev


func ConnectMongo() *mongo.Client {


	// get the utl string from the enviornment variable

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set")
	}

	// Create client options for database
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second)

	// Create context for connection to manage  time errors
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB with client options and context
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	// Ping MongoDB to verify connection for network check and stability
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}

	log.Println("MongoDB connected successfully")


	// return the client once everything is done
	return client
}


// return the collection if required in the apis 

// it takes the name of the collection and returns the collection

func OpenCollection(name string, client *mongo.Client) *mongo.Collection {
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		log.Fatal("DATABASE_NAME not set")
	}

	return client.Database(dbName).Collection(name)
}