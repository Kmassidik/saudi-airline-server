package repository

import (
	"api-server/config"
	"api-server/models"
	"errors"
	"log"
)

// GetAllBranchOffices retrieves all branch offices with pagination
func GetAllBranchOffices(limit, offset int) ([]models.BranchOfficeResponse, error) {
	var branchOffices []models.BranchOfficeResponse

	rows, err := config.DB.Query("SELECT id, name, address, total_counter FROM branch_offices ORDER BY id ASC LIMIT $1 OFFSET $2", limit, offset)
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

	row := config.DB.QueryRow("SELECT id, name, address, total_counter FROM branch_offices WHERE id = $1", id)
	err := row.Scan(&branchOffice.ID, &branchOffice.Name, &branchOffice.Address, &branchOffice.TotalCounter)
	if err != nil {
		log.Println("Error querying branch office by ID:", err)
		return nil, errors.New("branch office not found")
	}

	return &branchOffice, nil
}

// CreateBranchOffices creates a new branch office
func CreateBranchOffices(branchOffice *models.BranchOfficeCreateRequest) error {

	_, err := config.DB.Exec("INSERT INTO branch_offices (name, address, total_counter) VALUES ($1, $2, $3)",
		branchOffice.Name, branchOffice.Address, branchOffice.TotalCounter)
	if err != nil {
		log.Println("Error creating branch office:", err)
		return err
	}
	return nil
}

// UpdateBranchOffices updates an existing branch office by ID
func UpdateBranchOffices(id uint, branchOffice *models.BranchOfficeCreateRequest) error {
	_, err := config.DB.Exec("UPDATE branch_offices SET name = $1, address = $2, total_counter = $3 WHERE id = $4",
		branchOffice.Name, branchOffice.Address, branchOffice.TotalCounter, id)
	if err != nil {
		log.Println("Error updating branch office:", err)
		return err
	}
	return nil
}

// DeleteBranchOffices deletes a branch office by ID
func DeleteBranchOffices(id uint) error {
	_, err := config.DB.Exec("DELETE FROM branch_offices WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting branch office:", err)
		return err
	}
	return nil
}
