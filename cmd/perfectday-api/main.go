package main

import (
	"log"
	"os"
	"perfect-day/internal/api/server"
	"perfect-day/pkg/config"

	"github.com/joho/godotenv"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load .env file if it exists
	godotenv.Load()

	// Create configuration from environment variables
	cfg := &config.Config{
		DataDir:            getEnvOrDefault("DATA_DIR", "./.perfect-day"),
		GooglePlacesAPIKey: os.Getenv("GOOGLE_PLACES_API_KEY"),
	}

	// Create and start server
	srv := server.NewServer(cfg)
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}