package repository

import (
	"api-server/config"
	"api-server/helpers"
	"api-server/models"
	"api-server/repository/validation"
	"database/sql"
	"errors"
	"log"
)

// GetAllUsers retrieves all users from the database with pagination
func GetAllUsers(limit, offset int) ([]models.AllUserResponse, error) {
	var users []models.AllUserResponse

	rows, err := config.DB.Query(
		"SELECT id, full_name, email, role, likes, dislikes, image FROM users LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.AllUserResponse

		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Likes, &user.Dislikes, &user.Image); err != nil {
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

	// Query to retrieve the user by ID
	row := config.DB.QueryRow("SELECT id, full_name, email, role, likes, dislikes, image FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Likes, &user.Dislikes, &user.Image) // Scan the image name into imageName

	// If no rows are found
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with this ID
		}
		log.Println("Error scanning user:", err)
		return nil, err
	}

	// Return the user object directly (not wrapped in another object)
	return &user, nil
}

// CreateUser inserts a new user into the database with an optional image path
func CreateUser(user *models.User) error {
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

	// Hash the password if it's not empty (meaning it's being updated)
	if user.Password != "" {
		hashedPassword, err := helpers.HashingPasswordFunc(user.Password)
		if err != nil {
			log.Println("Error hashing password:", err)
			return errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword) // Store the hashed password
	}

	// Prepare the SQL query
	var query string
	if user.Image != "" {
		query = "UPDATE users SET full_name = $1, email = $2, password = $3, role = $4, likes = $5, dislikes = $6, image = $7 WHERE id = $8"
		_, err := config.DB.Exec(query, user.FullName, user.Email, user.Password, user.Role, user.Likes, user.Dislikes, user.Image, id)
		if err != nil {
			log.Println("Error updating user:", err)
			return err
		}
	} else {
		query = "UPDATE users SET full_name = $1, email = $2, password = $3, role = $4, likes = $5, dislikes = $6 WHERE id = $7"
		_, err := config.DB.Exec(query, user.FullName, user.Email, user.Password, user.Role, user.Likes, user.Dislikes, id)
		if err != nil {
			log.Println("Error updating user:", err)
			return err
		}
	}

	// No need to query the database again to get the existing image
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
