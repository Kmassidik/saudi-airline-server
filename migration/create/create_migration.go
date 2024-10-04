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

	// Check if the database exists
	if !databaseExists(db, dbName) {
		fmt.Printf("Database %s does not exist. Creating...\n", dbName)
		if err := createDatabase(db, dbName); err != nil {
			log.Fatalf("Failed to create database: %v\n", err)
		}
	}

	// Connect to the new database
	connStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s", dbUser, dbPassword, dbName, dbSslMode, dbHost)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define the SQL migration script
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS branch_offices (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		address TEXT,
		total_counter INT DEFAULT 0,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		image VARCHAR(255),
		role VARCHAR(50),
		likes INT DEFAULT 0,
		dislikes INT DEFAULT 0,
		branch_id INT,
		FOREIGN KEY (branch_id) REFERENCES branch_offices(id) ON DELETE CASCADE ON UPDATE CASCADE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS branch_counters (
		id SERIAL PRIMARY KEY,
		counter_location VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		branch_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (branch_id) REFERENCES branch_offices(id) ON DELETE CASCADE ON UPDATE CASCADE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS company_profiles (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		logo TEXT,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_feedback_history (
		id SERIAL PRIMARY KEY,
		likes INT DEFAULT 0,
		dislikes INT DEFAULT 0,
		officer_name VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create function to automatically update 'updatedAt' timestamp
	CREATE OR REPLACE FUNCTION update_timestamp_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updatedAt = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- Create triggers to update 'updatedAt' on row update for all tables
	CREATE TRIGGER update_users_updatedAt
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE FUNCTION update_timestamp_column();

	CREATE TRIGGER update_branch_offices_updatedAt
	BEFORE UPDATE ON branch_offices
	FOR EACH ROW
	EXECUTE FUNCTION update_timestamp_column();

	CREATE TRIGGER update_branch_counters_updatedAt
	BEFORE UPDATE ON branch_counters
	FOR EACH ROW
	EXECUTE FUNCTION update_timestamp_column();

	CREATE TRIGGER update_company_profiles_updatedAt
	BEFORE UPDATE ON company_profiles
	FOR EACH ROW
	EXECUTE FUNCTION update_timestamp_column();

	CREATE TRIGGER update_user_feedback_history_updatedAt
	BEFORE UPDATE ON user_feedback_history
	FOR EACH ROW
	EXECUTE FUNCTION update_timestamp_column();
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

// Create Database creates a new database
func createDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}
	return nil
}
