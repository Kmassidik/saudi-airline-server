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

	// Seed the database with branch offices and officers
	err = seedDatabase(db)
	if err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}

// seedDatabase inserts 10 branch offices and 50 officers for each branch office
func seedDatabase(db *sql.DB) error {
	// Insert 10 branch offices
	for i := 1; i <= 10; i++ {
		branchOfficeName := fmt.Sprintf("Branch Office %d", i)
		_, err := db.Exec("INSERT INTO branch_offices (name, address, total_counter) VALUES ($1, $2, $3)", branchOfficeName, fmt.Sprintf("Address of %s", branchOfficeName), 10)
		if err != nil {
			return fmt.Errorf("failed to insert branch office %d: %w", i, err)
		}

		// Insert 50 officers for each branch office
		for j := 1; j <= 25; j++ {
			officerName := fmt.Sprintf("Officer %d - %s", j, branchOfficeName)
			email := fmt.Sprintf("officer%d_%d@example.com", j, i)
			password := "password123" // You can use hashed password in real applications
			image := "image.png"
			role := "officer"

			// Insert officer linked to the current branch office
			_, err := db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES ($1, $2, $3, $4, $5, $6)", officerName, email, password, image, role, i)
			if err != nil {
				return fmt.Errorf("failed to insert officer %d for branch office %d: %w", j, i, err)
			}
		}
	}

	return nil
}
