package db

import (
	"context"
	"fmt"
	"notes-api/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this function is responsible for establishing a connection to the MongoDB database. It takes a config.Config struct as an argument, which contains the necessary configuration details for connecting to the database (such as the MongoDB URI and database name). The function returns a pointer to a mongo.Client, a pointer to a mongo.Database, and an error if any occurs during the connection process.
func ConnectDB(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	// prevent the app from freezing if the database connection takes too long, we can set a timeout for the connection attempt. This is done using context.WithTimeout, which creates a context that will automatically cancel after a specified duration (in this case, 10 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Error pinging MongoDB: %v", err)
	}

	db := client.Database(cfg.MongoDB)

	fmt.Println("Connected to MongoDB!")
	return client, db, nil
}

func DisconnectDB(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("Error disconnecting from MongoDB: %v", err)
	}

	fmt.Println("Disconnected from MongoDB!")
	return nil
}
