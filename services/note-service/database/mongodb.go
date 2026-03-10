package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"note-service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(cfg config.MongoDBConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URL))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	for i := 1; i <= 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := client.Ping(ctx, nil)
		cancel()
		if err == nil {
			log.Println("Connected to MongoDB")
			return client, nil
		}
		log.Printf("Waiting for MongoDB... attempt %d/10", i)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to MongoDB after 10 retries")
}

func CloseMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}
