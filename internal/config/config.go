package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	MongoDB    string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	// godotenv.Load() reads the environment variables from the .env file and loads them into the process's environment.
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// os.Getenv() retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present.
	//  instead of hardcoding the values, we can use a helper function.
	// config := &Config{
	// 	MongoURI:     os.Getenv("MONGODB_URI"),
	// 	MongoDB: os.Getenv("DATABASE_NAME"),
	// 	ServerPort:        os.Getenv("PORT"),
	// }

	mongoURI, err := extractENV("MONGODB_URI")
	if err != nil {
		return nil, err
	}

	mongoDB, err := extractENV("DATABASE_NAME")
	if err != nil {
		return nil, err
	}

	serverPort, err := extractENV("PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		MongoURI:   mongoURI,
		MongoDB:    mongoDB,
		ServerPort: serverPort,
	}, nil
}

func extractENV(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return val, nil
}
