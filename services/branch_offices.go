package services

import (
	"api-server/models"
	"api-server/repository"
)

// GetAllBranchOffices retrieves all branch offices with pagination
func GetAllBranchOffices(limit, offset int) ([]models.BranchOfficeResponse, error) {
	return repository.GetAllBranchOffices(limit, offset)
}

func GetAllBranchOfficesOptionList() ([]models.BranchOfficeOptionResponse, error) {
	return repository.GetAllBranchOfficesOption()
}

// GetBranchOfficeByID retrieves a branch office by ID
func GetBranchOfficeByID(id uint) (*models.BranchOfficeResponse, error) {
	return repository.GetBranchOfficesById(id)
}

// CreateBranchOffice creates a new branch office
func CreateBranchOffice(branchOffice *models.BranchOfficeCreateRequest) error {
	if err := repository.CreateBranchOffice(branchOffice); err != nil {
		return err
	}

	return nil
}

// UpdateBranchOffice updates an existing branch office by ID
func UpdateBranchOffice(id uint, branchOffice *models.BranchOfficeCreateRequest) error {
	return repository.UpdateBranchOffices(id, branchOffice)
}

// DeleteBranchOffice deletes a branch office by ID
func DeleteBranchOffice(id uint) error {
	return repository.DeleteBranchOffices(id)
}
