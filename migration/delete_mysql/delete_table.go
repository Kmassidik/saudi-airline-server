package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Use MySQL driver
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

	// Connect to the MySQL server
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/", dbUser, dbPassword, dbHost) // Connect to MySQL without a specific database
	db, err := sql.Open("mysql", connStr)
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
	query := fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", dbName)
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
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)) // Use backticks for database name
	if err != nil {
		return fmt.Errorf("error dropping database: %w", err)
	}
	return nil
}
