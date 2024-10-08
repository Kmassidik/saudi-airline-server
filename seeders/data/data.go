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

	errTotalData := insertTotalData(db)
	if errTotalData != nil {
		return fmt.Errorf("failed to insert total data: %w", errTotalData)
	}

	// Insert 1 Administrator
	err := insertAdministrator(db)
	if err != nil {
		return fmt.Errorf("failed to insert administrator: %w", err)
	}

	for i := 1; i <= 10; i++ {
		branchOfficeName := fmt.Sprintf("Branch Office %d", i)
		err := insertBranchOffice(db, branchOfficeName, i)
		if err != nil {
			return fmt.Errorf("failed to insert branch office %d: %w", i, err)
		}

		// Insert 1 Admin for each branch office
		err = insertAdmin(db, i, branchOfficeName)
		if err != nil {
			return fmt.Errorf("failed to insert admin for branch office %d: %w", i, err)
		}

		// Insert 1 Supervisor for each branch office
		err = insertSupervisor(db, i, branchOfficeName)
		if err != nil {
			return fmt.Errorf("failed to insert supervisor for branch office %d: %w", i, err)
		}

		// Insert 50 officers for each branch office
		for j := 1; j <= 50; j++ {
			err := insertOfficer(db, i, j, branchOfficeName)
			if err != nil {
				return fmt.Errorf("failed to insert officer %d for branch office %d: %w", j, i, err)
			}
		}

		// Insert total branch office data
		err = insertTotalBranchOffices(db, i, branchOfficeName)
		if err != nil {
			return fmt.Errorf("failed to insert total branch office data for %d: %w", i, err)
		}
	}

	return nil
}

// insertOfficer inserts an officer for a specific branch office
func insertOfficer(db *sql.DB, branchID, officerNum int, branchOfficeName string) error {
	officerName := fmt.Sprintf("Officer %d - %s", officerNum, branchOfficeName)
	email := fmt.Sprintf("officer%d_%d@example.com", officerNum, branchID)
	password, err := helpers.HashingPasswordFunc("password123")
	if err != nil {
		return err
	}
	image := "image.png"
	role := "officer"

	_, err = db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES ($1, $2, $3, $4, $5, $6)", officerName, email, password, image, role, branchID)
	return err
}

// insertAdministrator inserts an administrator for the main branch
func insertAdministrator(db *sql.DB) error {
	adminName := "Administrator"
	email := "administrator@example.com"
	password, err := helpers.HashingPasswordFunc("123456")
	if err != nil {
		return err
	}
	image := "admin.png"
	role := "administrator"

	_, err = db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES ($1, $2, $3, $4, $5, $6)", adminName, email, password, image, role, nil)
	return err
}

// insertAdmin inserts an admin for a specific branch office
func insertAdmin(db *sql.DB, branchID int, branchOfficeName string) error {
	adminName := fmt.Sprintf("Admin - %s", branchOfficeName)
	email := fmt.Sprintf("admin_%d@example.com", branchID)
	password, err := helpers.HashingPasswordFunc("password123")
	if err != nil {
		return err
	}
	image := "admin.png"
	role := "admin"

	_, err = db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES ($1, $2, $3, $4, $5, $6)", adminName, email, password, image, role, branchID)
	return err
}

// insertSupervisor inserts a supervisor for a specific branch office
func insertSupervisor(db *sql.DB, branchID int, branchOfficeName string) error {
	supervisorName := fmt.Sprintf("Supervisor - %s", branchOfficeName)
	email := fmt.Sprintf("supervisor_%d@example.com", branchID)
	password, err := helpers.HashingPasswordFunc("password123")

	if err != nil {
		return err
	}

	image := "supervisor.png"
	role := "supervisor"

	_, err = db.Exec("INSERT INTO users (full_name, email, password, image, role, branch_id) VALUES ($1, $2, $3, $4, $5, $6)", supervisorName, email, password, image, role, branchID)
	return err
}

// insertBranchOffice inserts a branch office into the database
func insertBranchOffice(db *sql.DB, name string, totalCounters int) error {
	_, err := db.Exec("INSERT INTO branch_offices (name, address, total_counter) VALUES ($1, $2, $3)", name, fmt.Sprintf("Address of %s", name), totalCounters)
	return err
}

// insertTotalBranchOffices inserts total data for branch offices
func insertTotalBranchOffices(db *sql.DB, branchID int, officeName string) error {
	_, err := db.Exec("INSERT INTO total_data_branch (name_office, total_likes, total_dislikes, branch_id) VALUES ($1, $2, $3, $4)", officeName, 0, 0, branchID)
	return err
}

// insertTotalBranchOffices inserts total data for branch offices
func insertTotalData(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO total_data (total_likes,total_dislikes,total_officer,total_voted) VALUES ($1, $2, $3, $4)", 0, 0, 0, 0)
	return err
}
