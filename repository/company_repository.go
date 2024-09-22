package repository

import (
	"api-server/config"
	"api-server/models"
	"database/sql"
	"log"
)

// Assuming you have a DB connection setup
var db *sql.DB

// SetDB sets the database connection for use in the repository
func SetDB(database *sql.DB) {
	db = database
}

// GetCompanyProfile fetches the company profile from the database
// GetCompanyProfile retrieves a company profile by ID
func GetCompanyProfile() (*models.CompanyProfile, error) {
	var company models.CompanyProfile

	// Query to select the company profile with a specific ID (1 in this case)
	row := config.DB.QueryRow(`SELECT name, logo FROM company_profiles WHERE id = $1`, 1)
	err := row.Scan(&company.Name, &company.Logo)

	// Check for errors during the scan
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are found, return nil and no error
			return nil, nil
		}
		// Log any other errors encountered during scanning
		log.Println("Error scanning company profile:", err)
		return nil, err // Return the error
	}

	// Return the company profile
	return &company, nil
}

// UpdateCompanyProfile updates the company profile in the database
func UpdateCompanyProfile(id uint, name string, logo string) error {
	query := `UPDATE company_profiles SET name = $1, logo = $2 WHERE id = $3`
	_, err := db.Exec(query, name, logo, id)
	if err != nil {
		log.Println("Error updating company profile:", err)
		return err
	}

	return nil
}
