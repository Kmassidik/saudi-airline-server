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

// GetBranchCountersByBranchID retrieves branch counters by branch ID, including names from related tables
func GetBranchCountersByBranchID(id uint) ([]models.BranchCounterWithNames, error) {
	query := `
        SELECT 
            bc.id, 
            bc.counter_location, 
            u.full_name AS full_name, 
            u.image AS image 
        FROM branch_counters bc
        JOIN branch_offices bo ON bc.branch_id = bo.id
        JOIN users u ON bc.user_id = u.id
        WHERE bc.branch_id = $1
    `

	rows, err := config.DB.Query(query, id) // Use Query instead of Exec for SELECT statements
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counters []models.BranchCounterWithNames
	for rows.Next() {
		var counter models.BranchCounterWithNames
		if err := rows.Scan(
			&counter.ID,
			&counter.CounterLocation,
			&counter.FullName,
			&counter.Image,
		); err != nil {
			return nil, err
		}

		counters = append(counters, counter)
	}

	return counters, nil
}

// DeleteBranchCounter deletes a branch counter by ID
func DeleteBranchCounter(id uint) error {
	_, err := config.DB.Exec("DELETE FROM branch_counters WHERE id = $1", id)
	return err
}
