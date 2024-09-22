package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hratsch/zesty-sips-api/internal/api"
	"github.com/hratsch/zesty-sips-api/internal/config"
	"github.com/hratsch/zesty-sips-api/internal/db"
	"github.com/joho/godotenv"
)

func connectToDatabase() (*sql.DB, error) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbName)

	return db.Connect(dbURL)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}


	// Initialize database connection
	database, err := connectToDatabase()
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
