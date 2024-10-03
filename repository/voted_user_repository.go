package repository

import (
	"api-server/config"
	"api-server/models"
	"fmt"
)

func VotedUserLike(voteType string, data *models.User) error {
	// Get the user ID from data
	userId := data.ID
	officerName := data.FullName

	// Prepare the query to update the users table
	var updateQuery string
	if voteType == "like" {
		updateQuery = `UPDATE users SET likes = likes + 1 WHERE id = $1`
	} else if voteType == "dislike" {
		updateQuery = `UPDATE users SET dislikes = dislikes + 1 WHERE id = $1`
	} else {
		return fmt.Errorf("invalid vote type")
	}

	// Execute the update query
	if _, err := config.DB.Exec(updateQuery, userId); err != nil {
		return err
	}

	// Prepare the insert query for user_feedback_history
	insertQuery := `INSERT INTO user_feedback_history (likes, dislikes, officer_name, user_id) VALUES ($1, $2, $3, $4)`
	likes := 0
	dislikes := 0

	if voteType == "like" {
		likes = 1
	} else if voteType == "dislike" {
		dislikes = 1
	}

	// Insert feedback history
	if _, err := config.DB.Exec(insertQuery, likes, dislikes, officerName, userId); err != nil {
		return err
	}

	return nil
}
