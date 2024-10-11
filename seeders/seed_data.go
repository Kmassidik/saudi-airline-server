package main

import (
	"api-server/helpers"
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

// seedDatabase inserts branch offices and associated users
func seedDatabase(db *sql.DB) error {

	// Insert 1 Administrator
	err := insertAdministrator(db)
	if err != nil {
		return fmt.Errorf("failed to insert administrator: %w", err)
	}

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

	return nil
}

// insertAdministrator inserts an administrator for the main branch
func insertAdministrator(db *sql.DB) error {
	adminName := "Administrator"
	email := "administrator@example.com"
	password, err := helpers.HashingPasswordFunc("admin12345")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	image := "admin.png"
	role := "administrator"

	// Insert administrator, ensuring branch_id is handled properly
	_, err = db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES ($1, $2, $3, $4, $5, $6)", adminName, email, password, image, role, nil)
	if err != nil {
		return fmt.Errorf("failed to insert administrator: %w", err)
	}
	return nil
}
