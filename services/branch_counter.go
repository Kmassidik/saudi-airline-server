package services

import (
	"api-server/models"
	"api-server/repository"
	"strconv"
)

// GetBranchCountersByBranchID retrieves branch counters for a specific branch ID
func GetBranchCountersByBranchID(id uint) ([]models.BranchCounterWithNames, error) {
	// Convert id from int to uint
	uintID := uint(id)
	return repository.GetBranchCountersByBranchID(uintID)
}

// CreateBranchCounter creates a new branch counter
func CreateBranchCounter(branchCounter *models.BranchCounter) error {
	return repository.CreateBranchCounter(branchCounter)
}

// DeleteBranchCounter deletes a branch counter by ID
func DeleteBranchCounter(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return repository.DeleteBranchCounter(uint(idInt))
}
