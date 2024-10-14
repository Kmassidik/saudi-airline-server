package repository

import (
	"api-server/config"
	"api-server/models"
	"fmt"
	"log"
)

func VotedUserLike(voteType string, data *models.User) error {
	// Get the user ID from data
	userId := data.ID
	branchId := data.BranchId
	officerName := data.FullName

	// Begin a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // Rollback on error
		}
	}()

	// Prepare the query to update the users table
	var updateQuery string
	if voteType == "like" {
		updateQuery = `UPDATE users SET likes = likes + 1 WHERE id = ?`
	} else if voteType == "dislike" {
		updateQuery = `UPDATE users SET dislikes = dislikes + 1 WHERE id = ?`
	} else {
		return fmt.Errorf("invalid vote type")
	}

	// Execute the update query
	if _, err := tx.Exec(updateQuery, userId); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Prepare the insert query for user_feedback_history
	insertQuery := `INSERT INTO user_feedback_history (likes, dislikes, officer_name, user_id, branch_id) VALUES (?, ?, ?, ?, ?)`
	likes := 0
	dislikes := 0

	if voteType == "like" {
		likes = 1
	} else if voteType == "dislike" {
		dislikes = 1
	}

	// Insert feedback history
	if _, err := tx.Exec(insertQuery, likes, dislikes, officerName, userId, branchId); err != nil {
		return fmt.Errorf("failed to insert user feedback history: %v", err)
	}

	// Commit the transaction if no errors occurred
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
