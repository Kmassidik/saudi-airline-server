package repository

import (
	"api-server/config"
	"api-server/models"
	"fmt"
	"log"
)

func TotalDataDashboard() (int, int, int, int, error) {
	var totalOfficer, totalLikes, totalDislikes, totalVoted int

	// Combine all counts into a single query
	query := `SELECT total_officer, total_likes, total_dislikes, total_voted FROM total_data WHERE id = 1;`

	// Execute the query
	row := config.DB.QueryRow(query)

	// Scan the results into the respective variables
	err := row.Scan(&totalOfficer, &totalLikes, &totalDislikes, &totalVoted)
	if err != nil {
		log.Println("Error querying total data dashboard:", err)
		return 0, 0, 0, 0, err
	}

	return totalOfficer, totalLikes, totalDislikes, totalVoted, nil
}

func TotalDataBranchDashboard() ([]models.BranchData, error) {
	rows, err := config.DB.Query("SELECT id, name_office, total_likes, total_dislikes, branch_id FROM total_data_branch ORDER BY total_likes DESC")
	if err != nil {
		log.Println("Error querying total data dashboard:", err)
		return nil, err
	}
	defer rows.Close()

	var branchDataList []models.BranchData

	for rows.Next() {
		var branchData models.BranchData
		err := rows.Scan(&branchData.ID, &branchData.NameOffice, &branchData.TotalLikes, &branchData.TotalDislikes, &branchData.BranchID)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		branchDataList = append(branchDataList, branchData)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		return nil, err
	}

	return branchDataList, nil
}

func DataOfficerDashboard(limit uint, offset uint) ([]models.DashboardUsers, error) {
	// Prepare the SQL query to select the desired fields
	query := "SELECT full_name, likes, dislikes FROM users WHERE role = 'officer' ORDER BY likes DESC LIMIT $1 OFFSET $2"
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.DashboardUsers

	// Loop through the rows to populate the users slice
	for rows.Next() {
		var user models.DashboardUsers
		if err := rows.Scan(&user.Name, &user.Likes, &user.Dislikes); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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
