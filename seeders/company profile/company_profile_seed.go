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
	dbHost := os.Getenv("DB_HOST")

	// Define the connection string for PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s", dbUser, dbPassword, dbName, dbSslMode, dbHost)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if data already exists in the `company` table
	row := db.QueryRow("SELECT COUNT(*) FROM company_profiles")
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Fatalf("Error counting rows: %v\n", err)
	}

	// If no data exists, insert seed data
	if count == 0 {
		seedSQL := `
		INSERT INTO company_profiles (name, logo) 
		VALUES ('Sample Company', 'application_logo.png');
		`
		_, err = db.Exec(seedSQL)
		if err != nil {
			log.Fatalf("Failed to insert seed data: %v\n", err)
		}
		fmt.Println("Seed data inserted successfully.")
	} else {
		fmt.Println("Data already exists, skipping seed.")
	}
}
