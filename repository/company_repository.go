package repository

import (
	"api-server/config"
	"api-server/models"
	"database/sql"
	"log"
)

// GetCompanyProfile fetches the company profile from the database
func GetCompanyProfile() (*models.CompanyProfile, error) {
	var company models.CompanyProfile

	// Query to select the company profile with a specific ID (1 in this case)
	row := config.DB.QueryRow(`SELECT name, logo FROM company_profiles WHERE id = ?`, 1)
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
	// Prepare the SQL query
	var query string

	if logo != "" {
		query = "UPDATE company_profiles SET name = ?, logo = ? WHERE id = ?"
		_, err := config.DB.Exec(query, name, logo, id)
		if err != nil {
			log.Println("Error updating company profile:", err)
			return err
		}
	} else {
		query = "UPDATE company_profiles SET name = ? WHERE id = ?"
		_, err := config.DB.Exec(query, name, id)
		if err != nil {
			log.Println("Error updating company profile:", err)
			return err
		}
	}

	return nil
}
