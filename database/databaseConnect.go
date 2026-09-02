package database

import(
	"context"
	"fmt"
	"log"
	"time"
	"os"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"github.com/joho/godotenv"
	
)
func DatabaseConnect() *mongo.Client {
	fmt.Println("Database Connecting...")

	// Only try to load .env in development
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("No .env file found - using system environment variables")
		}
	}

	MongoDB := os.Getenv("MONGO_DB_URL_STRING")
	fmt.Println("MongoDB URL:", MongoDB) // Be careful - remove this in production for security!

	if MongoDB == "" {
		log.Fatal("MONGO_DB_URL_STRING environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}
	
	fmt.Println("MongoDB Connected Successfully")
	return client
}

// func DatabaseConnect() *mongo.Client{
// 	fmt.Println("Database Connected")

// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Printf("Error loading .env file: %v", err)
// 	}

// 	MongoDB := os.Getenv("MONGO_DB_URL_STRING")
// 	fmt.Println("MongoDB URL:", MongoDB)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(options.Client().ApplyURI(MongoDB))
// 	if err != nil {
// 		log.Fatal("Error connecting to MongoDB: ", err)
// 	}

// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatal("Could not ping MongoDB: ", err)
// 	}
// 	fmt.Println("MongoDB Connected Successfully")
// 	return client
	
// }

var Client *mongo.Client = DatabaseConnect()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection
}