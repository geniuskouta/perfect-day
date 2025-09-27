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
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Debug: Print environment variables
	log.Printf("FRONTEND_URL from env: %s", os.Getenv("FRONTEND_URL"))
	log.Printf("GOOGLE_PLACES_API_KEY set: %t", os.Getenv("GOOGLE_PLACES_API_KEY") != "")

	// Create configuration from environment variables
	cfg := &config.Config{
		DataDir:            getEnvOrDefault("DATA_DIR", "./.perfect-day"),
		GooglePlacesAPIKey: os.Getenv("GOOGLE_PLACES_API_KEY"),
		FrontendUrl:        getEnvOrDefault("FRONTEND_URL", "http://localhost:5173"),
	}

	log.Printf("Final config FrontendUrl: %s", cfg.FrontendUrl)

	// Create and start server
	srv := server.NewServer(cfg)
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}