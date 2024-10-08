package repository

import (
	"api-server/config"
	"fmt"
	"log"
)

func TotalDataDashboard() (int, int, int, int, error) {
	var totalUsers, totalLikes, totalDislikes, totalVoted int

	// Combine all counts into a single query
	query := `
		SELECT 
			COUNT(id) AS total_users, 
			SUM(likes) AS total_likes, 
			SUM(dislikes) AS total_dislikes, 
			(SELECT COUNT(id) FROM user_feedback_history) AS total_voted 
		FROM users 
		WHERE role = 'officer';`

	// Execute the query
	row := config.DB.QueryRow(query)

	// Scan the results into the respective variables
	err := row.Scan(&totalUsers, &totalLikes, &totalDislikes, &totalVoted)
	if err != nil {
		log.Println("Error querying total data dashboard:", err)
		return 0, 0, 0, 0, err
	}

	return totalUsers, totalLikes, totalDislikes, totalVoted, nil
}

func TotalDataBranchDashboard(id uint, option string) error {

	if option == "week" {
	} else if option == "month" {

	}
	return nil
}

func DiagramDataDashboard() {

}

func UpdateDashboard(branchId uint, voteType string) error {
	var (
		totalUpdateQuery  string
		branchUpdateQuery string
	)

	// Determine the update queries based on vote type
	switch voteType {
	case "like":
		totalUpdateQuery = `UPDATE total_data SET total_likes = total_likes + 1, total_voted = total_voted + 1 WHERE id = 1`
		branchUpdateQuery = `UPDATE total_data_branch SET total_likes = total_likes + 1 WHERE branch_id = $1`
	case "dislike":
		totalUpdateQuery = `UPDATE total_data SET total_dislikes = total_dislikes + 1, total_voted = total_voted + 1 WHERE id = 1`
		branchUpdateQuery = `UPDATE total_data_branch SET total_dislikes = total_dislikes + 1 WHERE branch_id = $1`
	default:
		return fmt.Errorf("invalid vote type")
	}

	// Execute total_data update
	if _, err := config.DB.Exec(totalUpdateQuery); err != nil {
		return fmt.Errorf("failed to update total_data: %v", err)
	}

	// Execute total_data_branch update
	if _, err := config.DB.Exec(branchUpdateQuery, branchId); err != nil {
		return fmt.Errorf("failed to update total_data_branch: %v", err)
	}

	return nil
}
