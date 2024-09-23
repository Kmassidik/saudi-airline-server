package services

import (
	"api-server/models"
	"api-server/repository"
)

// CreateBranchCounter creates a new branch counter
func CreateBranchCounter(branchCounter *models.BranchCounter) error {
	return repository.CreateBranchCounter(branchCounter)
}

// UpdateBranchCounter updates an existing branch counter
func UpdateBranchCounter(id uint, branchCounter *models.BranchCounter) error {
	return repository.UpdateBranchCounter(id, branchCounter)
}

// DeleteBranchCounter deletes a branch counter by ID
func DeleteBranchCounter(id uint) error {
	return repository.DeleteBranchCounter(id)
}
