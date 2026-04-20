package main

//  work flow is
// 1. Load configuration from .env file using config.LoadConfig().
// 2. Connect to the MongoDB database using db.ConnectDB() with the loaded configuration.
// 3. Create a new router using server.NewRouter().
// 4. Start the HTTP server on the specified port using router.Run().

import (
	"fmt"
	"log"
	"notes-api/internal/config"
	"notes-api/internal/db"
	"notes-api/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, database, err := db.ConnectDB(*cfg)

	if err != nil {
		log.Fatalf("Failed to connect to database:: %v", err)
	}
	// if any issueocurs we are disconnecting the database connection and logging the error if it occurs.
	defer func() {
		if err := db.DisconnectDB(client); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	router := server.NewRouter(database)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
