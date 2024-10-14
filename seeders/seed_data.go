package main

import (
	"api-server/helpers"
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

	// Define the connection string for MySQL
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Seed the database with branch offices, officers, and total data
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

	// Check if data already exists in the `company_profiles` table
	row := db.QueryRow("SELECT COUNT(*) FROM company_profiles")
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Fatalf("Error counting rows: %v\n", err)
	}

	// If no data exists, insert seed data for company_profiles
	if count == 0 {
		seedSQL := `
		INSERT INTO company_profiles (name, logo) 
		VALUES ('Sample Company', 'application_logo.png');
		`
		_, err = db.Exec(seedSQL)
		if err != nil {
			log.Fatalf("Failed to insert seed data: %v\n", err)
		}
		fmt.Println("Seed data for company_profiles inserted successfully.")
	} else {
		fmt.Println("Data already exists in company_profiles, skipping seed.")
	}

	// Seed total_data
	err = seedTotalData(db)
	if err != nil {
		return fmt.Errorf("failed to seed total_data: %w", err)
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
	_, err = db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES (?, ?, ?, ?, ?, ?)", adminName, email, password, image, role, nil)
	if err != nil {
		return fmt.Errorf("failed to insert administrator: %w", err)
	}
	return nil
}

// seedTotalData inserts seed data into the total_data table
func seedTotalData(db *sql.DB) error {
	// Check if data already exists in the total_data table
	row := db.QueryRow("SELECT COUNT(*) FROM total_data")
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatalf("Error counting rows in total_data: %v\n", err)
	}

	// If no data exists, insert seed data
	if count == 0 {
		seedSQL := `
		INSERT INTO total_data (total_likes, total_dislikes, total_officer, total_voted) 
		VALUES (0, 0, 0, 0); 
		`
		_, err = db.Exec(seedSQL)
		if err != nil {
			return fmt.Errorf("failed to insert seed data into total_data: %w", err)
		}
		fmt.Println("Seed data for total_data inserted successfully.")
	} else {
		fmt.Println("Data already exists in total_data, skipping seed.")
	}

	return nil
}
