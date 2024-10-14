package repository

import (
	"api-server/config"
	"api-server/models"
	"errors"
	"fmt"
	"log"
)

// GetAllBranchOffices retrieves all branch offices with pagination
func GetAllBranchOffices(limit, offset int) ([]models.BranchOfficeResponse, error) {
	var branchOffices []models.BranchOfficeResponse

	// MySQL uses ? as a placeholder
	rows, err := config.DB.Query("SELECT id, name, address, total_counter FROM branch_offices ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("Error querying branch offices:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branchOffice models.BranchOfficeResponse
		if err := rows.Scan(&branchOffice.ID, &branchOffice.Name, &branchOffice.Address, &branchOffice.TotalCounter); err != nil {
			log.Println("Error scanning branch office:", err)
			return nil, err
		}
		branchOffices = append(branchOffices, branchOffice)
	}

	return branchOffices, nil
}

// GetAllBranchOfficesOption retrieves all branch offices for options
func GetAllBranchOfficesOption() ([]models.BranchOfficeOptionResponse, error) {
	var branchOffices []models.BranchOfficeOptionResponse

	// MySQL uses ? as a placeholder
	rows, err := config.DB.Query("SELECT id, name FROM branch_offices ORDER BY id ASC")
	if err != nil {
		log.Println("Error querying branch offices:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branchOffice models.BranchOfficeOptionResponse
		if err := rows.Scan(&branchOffice.ID, &branchOffice.Name); err != nil {
			log.Println("Error scanning branch office:", err)
			return nil, err
		}
		branchOffices = append(branchOffices, branchOffice)
	}

	return branchOffices, nil
}

// GetBranchOfficesCount retrieves the total number of branch offices
func GetBranchOfficesCount() (int, error) {
	var count int
	row := config.DB.QueryRow("SELECT COUNT(*) FROM branch_offices")
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error querying branch offices count:", err)
		return 0, err
	}
	return count, nil
}

// GetBranchOfficesById retrieves a branch office by its ID
func GetBranchOfficesById(id uint) (*models.BranchOfficeResponse, error) {
	var branchOffice models.BranchOfficeResponse

	row := config.DB.QueryRow("SELECT id, name, address, total_counter FROM branch_offices WHERE id = ?", id)
	err := row.Scan(&branchOffice.ID, &branchOffice.Name, &branchOffice.Address, &branchOffice.TotalCounter)
	if err != nil {
		return nil, errors.New("branch office not found")
	}

	return &branchOffice, nil
}

// CreateBranchOffice creates a new branch office and returns its ID.
func CreateBranchOffice(branchOffice *models.BranchOfficeCreateRequest) error {
	// Begin a new transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback in case of an error

	// Insert the branch office
	_, err = tx.Exec(
		"INSERT INTO branch_offices (name, address, total_counter) VALUES (?, ?, ?)",
		branchOffice.Name, branchOffice.Address, branchOffice.TotalCounter,
	)
	if err != nil {
		log.Println("Error creating branch office:", err)
		return err
	}

	// Get the last inserted ID
	var branchID int
	err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&branchID)
	if err != nil {
		log.Println("Error getting last inserted ID:", err)
		return err
	}

	// Insert into total_data_branch using the generated branch ID
	_, err = tx.Exec(
		"INSERT INTO total_data_branch (name_office, total_likes, total_dislikes, branch_id) VALUES (?, ?, ?, ?)",
		branchOffice.Name, 0, 0, branchID,
	)
	if err != nil {
		log.Println("Error creating total data for branch office:", err)
		return err
	}

	// Commit the transaction if everything is successful
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil // Return the branch ID on success
}

// UpdateBranchOffices updates an existing branch office by ID
func UpdateBranchOffices(id uint, branchOffice *models.BranchOfficeCreateRequest) error {
	result, err := config.DB.Exec("UPDATE branch_offices SET name = ?, address = ?, total_counter = ? WHERE id = ?",
		branchOffice.Name, branchOffice.Address, branchOffice.TotalCounter, id)
	if err != nil {
		log.Println("Error updating branch office:", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no branch office found with the given ID: %d", id)
	}

	return nil
}

// DeleteBranchOffices deletes a branch office by ID
func DeleteBranchOffices(id uint) error {
	result, err := config.DB.Exec("DELETE FROM branch_offices WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting branch office:", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no branch office found with the given ID: %d", id)
	}

	return nil
}
