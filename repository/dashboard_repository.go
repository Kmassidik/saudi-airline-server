package repository

import (
	"api-server/config"
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
