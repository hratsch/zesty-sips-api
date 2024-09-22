package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hratsch/zesty-sips-api/internal/api"
	"github.com/hratsch/zesty-sips-api/internal/config"
	"github.com/hratsch/zesty-sips-api/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize router
	router := api.NewRouter(database)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
