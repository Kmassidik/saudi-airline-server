package repository

import (
	"api-server/config"
	"api-server/models"
	"log"
)

// CreateBranchCounter creates a new branch counter
func CreateBranchCounter(branchCounter *models.BranchCounter) error {
	_, err := config.DB.Exec("INSERT INTO branch_counters (counter_location, user_id, branch_id) VALUES ($1, $2, $3)",
		branchCounter.CounterLocation, branchCounter.UserID, branchCounter.BranchID)

	if err != nil {
		log.Println("Error creating branch counter:", err)
		return err
	}

	return nil
}

// UpdateBranchCounter updates an existing branch counter
func UpdateBranchCounter(id uint, branchCounter *models.BranchCounter) error {
	_, err := config.DB.Exec("UPDATE branch_counters SET counter_location = $1, user_id = $2, branch_id = $3 WHERE id = $4", branchCounter.CounterLocation, branchCounter.UserID, branchCounter.BranchID, id)
	return err
}

// DeleteBranchCounter deletes a branch counter by ID
func DeleteBranchCounter(id uint) error {
	_, err := config.DB.Exec("DELETE FROM branch_counters WHERE id = $1", id)
	return err
}
