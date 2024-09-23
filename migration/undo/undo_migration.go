package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")
	dbHost := os.Getenv("DB_HOST") // Host of the PostgreSQL server

	// Connect to default PostgreSQL database
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=%s host=%s", dbUser, dbPassword, dbSslMode, dbHost)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete the database if it exists
	if databaseExists(db, dbName) {
		fmt.Printf("Database %s exists. Deleting...\n", dbName)
		if err := deleteDatabase(db, dbName); err != nil {
			log.Fatalf("Failed to delete database: %v\n", err)
		}
		fmt.Printf("Database %s deleted successfully.\n", dbName)
	} else {
		fmt.Printf("Database %s does not exist.\n", dbName)
	}
}

// databaseExists checks if the database exists
func databaseExists(db *sql.DB, dbName string) bool {
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error checking if database exists: %v\n", err)
	}
	defer rows.Close()
	return rows.Next()
}

// deleteDatabase deletes an existing database
func deleteDatabase(db *sql.DB, dbName string) error {
	// Drop the database
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("error dropping database: %w", err)
	}
	return nil
}
