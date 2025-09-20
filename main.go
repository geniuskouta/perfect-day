package main

import (
	"log"
	"perfect-day/src/cli"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env file doesn't exist, just log it
		log.Printf("No .env file found or error loading it: %v", err)
	}

	cli.Execute()
}