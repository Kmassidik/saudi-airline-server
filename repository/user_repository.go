package repository

import (
	"api-server/config"
	"api-server/helpers"
	"api-server/models"
	"api-server/repository/validation"
	"database/sql"
	"log"
)

// GetAllUsers retrieves all users from the database with pagination
func GetAllUsers(limit, offset int) ([]models.User, error) {
	var users []models.User

	rows, err := config.DB.Query(
		"SELECT id, full_name, email, password, role, likes, dislikes, image FROM users LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Password, &user.Role, &user.Likes, &user.Dislikes, &user.Image); err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iteration:", err)
		return nil, err
	}

	return users, nil
}

// GetUsersCount retrieves the total number of users
func GetUsersCount() (int, error) {
	var count int

	row := config.DB.QueryRow("SELECT COUNT(*) FROM users")
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error querying users count:", err)
		return 0, err
	}

	return count, nil
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(id uint) (*models.User, error) {
	var user models.User

	row := config.DB.QueryRow("SELECT id, full_name, email, password, role, likes, dislikes FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Password, &user.Role, &user.Likes, &user.Dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with this ID
		}
		log.Println("Error scanning user:", err)
		return nil, err
	}

	return &user, nil
}

// CreateUser inserts a new user into the database with an optional image path
func CreateUser(user *models.User) error {
	// Perform validation before insertion
	if err := validation.ValidateUser(user); err != nil {
		return err
	}

	// Hash the user's password before saving to the database
	hashedPassword, err := helpers.HashingPasswordFunc(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	log.Println("Inserting user:", user)
	// Insert the user into the database, including the image path
	_, err = config.DB.Exec(
		"INSERT INTO users (full_name, email, password, role, likes, dislikes, image) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.FullName, user.Email, user.Password, user.Role, user.Likes, user.Dislikes, user.Image,
	)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}

	return nil
}

// UpdateUser updates an existing user by ID in the database
func UpdateUser(id uint, user *models.User) error {
	// Perform validation before updating
	if err := validation.ValidateUser(user); err != nil {
		return err
	}

	_, err := config.DB.Exec(
		"UPDATE users SET full_name = $1, email = $2, password = $3, role = $4, likes = $5, dislikes = $6 WHERE id = $7",
		user.FullName, user.Email, user.Password, user.Role, user.Likes, user.Dislikes, id,
	)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}

// DeleteUser deletes a user by ID from the database
func DeleteUser(id uint) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	return nil
}
