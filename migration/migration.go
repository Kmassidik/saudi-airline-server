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
	// Load environment variables from .env file located outside the migration folder
	err := godotenv.Load(".env") // Adjust the path according to your folder structure
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")
	dbHost := os.Getenv("DB_HOST") // Added host

	// Define the connection string for PostgreSQL (initially connecting to default database)
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=%s host=%s", dbUser, dbPassword, dbSslMode, dbHost)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the database exists
	if !databaseExists(db, dbName) {
		fmt.Printf("Database %s does not exist. Creating...\n", dbName)
		if err := createDatabase(db, dbName); err != nil {
			log.Fatalf("Failed to create database: %v\n", err)
		}
	}

	// Connect to the newly created database
	connStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s", dbUser, dbPassword, dbName, dbSslMode, dbHost)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define the migration SQL script
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		image VARCHAR(255),
		role VARCHAR(50),
		likes INT DEFAULT 0,
		dislikes INT DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS branch_offices (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		address TEXT,
		total_counter INT DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS branch_counters (
		id SERIAL PRIMARY KEY,
		counter_location VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		branch_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (branch_id) REFERENCES branch_offices(id) ON DELETE CASCADE ON UPDATE CASCADE
	);

	CREATE TABLE IF NOT EXISTS company_profiles (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		logo TEXT
	);
	`

	// Execute the migration script
	_, err = db.Exec(migrationSQL)
	if err != nil {
		log.Fatalf("Failed to execute migration: %v\n", err)
	}

	fmt.Println("Migration executed successfully!")
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

// createDatabase creates a new database
func createDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}
	return nil
}
