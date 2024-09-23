package services

import (
	"api-server/models"
	"api-server/repository"
)

// GetAllBranchOffices retrieves all branch offices with pagination
func GetAllBranchOffices(limit, offset int) ([]models.BranchOffice, error) {
	// Fetch branch offices from the repository
	return repository.GetAllBranchOffices(limit, offset)
}

// GetBranchOfficeByID retrieves a branch office by ID
func GetBranchOfficeByID(id uint) (*models.BranchOffice, error) {
	// Fetch a single branch office by ID from the repository
	return repository.GetBranchOfficesById(id)
}

// CreateBranchOffice creates a new branch office
func CreateBranchOffice(branchOffice *models.BranchOffice) error {
	// Perform any validation before inserting into the repository if necessary
	if err := repository.CreateBranchOffices(branchOffice); err != nil {
		return err
	}

	return nil
}

// UpdateBranchOffice updates an existing branch office by ID
func UpdateBranchOffice(id uint, branchOffice *models.BranchOffice) error {
	// Call the repository to update the branch office data
	return repository.UpdateBranchOffices(id, branchOffice)
}

// DeleteBranchOffice deletes a branch office by ID
func DeleteBranchOffice(id uint) error {
	// Call the repository to delete the branch office by ID
	return repository.DeleteBranchOffices(id)
}
