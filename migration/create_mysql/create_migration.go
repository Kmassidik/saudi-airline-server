package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"
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
	dbHost := os.Getenv("DB_HOST") // Host of the MySQL server

	// Connect to MySQL server (using a default database like mysql)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql", dbUser, dbPassword, dbHost)
	db, err := sql.Open("mysql", connStr)
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
	connStr = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define the SQL migration script as individual statements
	migrationSQLs := []string{
		`CREATE TABLE IF NOT EXISTS branch_offices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		address TEXT,
		total_counter INT DEFAULT 0,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
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
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,

		`CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);`,

		`CREATE TABLE IF NOT EXISTS branch_counters (
		id INT AUTO_INCREMENT PRIMARY KEY,
		counter_location VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		branch_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (branch_id) REFERENCES branch_offices(id) ON DELETE CASCADE ON UPDATE CASCADE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS company_profiles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		logo TEXT,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS user_feedback_history (
		id INT AUTO_INCREMENT PRIMARY KEY,
		likes INT DEFAULT 0,
		dislikes INT DEFAULT 0,
		officer_name VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		branch_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (branch_id) REFERENCES branch_offices(id) ON DELETE CASCADE ON UPDATE CASCADE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS total_data (
		id INT AUTO_INCREMENT PRIMARY KEY,
		total_likes INT DEFAULT 0,
		total_dislikes INT DEFAULT 0,
		total_officer INT DEFAULT 0,
		total_voted INT DEFAULT 0,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS total_data_branch (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name_office VARCHAR(255) NOT NULL,
		total_likes INT DEFAULT 0,
		total_dislikes INT DEFAULT 0,
		branch_id INT NOT NULL,
		FOREIGN KEY (branch_id) REFERENCES branch_offices(id) ON DELETE CASCADE ON UPDATE CASCADE,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`,
	}

	// Execute each migration statement
	for _, sql := range migrationSQLs {
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatalf("Failed to execute migration: %v\n", err)
		}
	}

	fmt.Println("Migration executed successfully!")

}

// databaseExists checks if the database exists
func databaseExists(db *sql.DB, dbName string) bool {
	query := fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", dbName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error checking if database exists: %v\n", err)
	}
	defer rows.Close()
	return rows.Next()
}

// createDatabase creates a new database
func createDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", dbName))
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}
	return nil
}
