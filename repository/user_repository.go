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
func GetAllUsers(limit, offset int, role string) ([]models.UserAllResponse, error) {
	var users []models.UserAllResponse
	var rows *sql.Rows
	var err error

	// Handle role-based query: "officer" or not "officer"
	if role == "officer" {
		rows, err = config.DB.Query(
			"SELECT id, full_name, email, role, likes, dislikes, image, branch_id FROM users WHERE role = $3 ORDER BY id ASC LIMIT $1 OFFSET $2",
			limit, offset, role,
		)
	} else {
		rows, err = config.DB.Query(
			"SELECT id, full_name, email, role, likes, dislikes, image, branch_id FROM users WHERE role != $3 ORDER BY id ASC LIMIT $1 OFFSET $2",
			limit, offset, "officer",
		)
	}

	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserAllResponse

		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Likes, &user.Dislikes, &user.Image, &user.BranchId); err != nil {
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
func GetUsersCount(role string) (int, error) {
	var count int
	var row *sql.Row

	if role == "officer" {
		// Use a parameterized query to safely query based on role
		row = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'officer'")
	} else {
		// When role is not provided, count all users
		row = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role != 'officer'")
	}

	// Scan the result into the count variable
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
	row := config.DB.QueryRow("SELECT id, full_name, email, role, likes, dislikes, image, branch_id, password  FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Likes, &user.Dislikes, &user.Image, &user.BranchId, &user.Password) // Scan the image name into imageName

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
	// Insert the user into the database, including the image path
	_, err = config.DB.Exec(
		"INSERT INTO users (full_name, email, password, role, likes, dislikes, image, branch_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		user.FullName, user.Email, user.Password, user.Role, user.Likes, user.Dislikes, user.Image, user.BranchId,
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
		query = "UPDATE users SET full_name = $1, email = $2, password = $3, role = $4, image = $5, branch_id = $6 WHERE id = $7"
		_, err := config.DB.Exec(query, user.FullName, user.Email, user.Password, user.Role, user.Image, user.BranchId, id)
		if err != nil {
			log.Println("Error updating user:", err)
			return err
		}
	} else {
		query = "UPDATE users SET full_name = $1, email = $2, password = $3, role = $4, branch_id = $5 WHERE id = $6"
		_, err := config.DB.Exec(query, user.FullName, user.Email, user.Password, user.Role, user.BranchId, id)
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

func GetAllUsersByBranchOfiice(branchId uint) ([]models.UserByBranchOfiiceResponse, error) {
	var users []models.UserByBranchOfiiceResponse

	// Execute a SELECT query to fetch users by branch_id
	rows, err := config.DB.Query("SELECT id, full_name FROM users WHERE branch_id = $1 AND role = 'officer' AND NOT EXISTS ( SELECT 1 FROM	branch_counters WHERE branch_counters.user_id = users.id);", branchId)
	if err != nil {
		log.Println("Error fetching users by branch:", err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and scan each user
	for rows.Next() {
		var user models.UserByBranchOfiiceResponse
		if err := rows.Scan(&user.ID, &user.FullName); err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration:", err)
		return nil, err
	}

	return users, nil
}

func CheckUserAuthentication(email string, password string) (models.User, error) {
	var user models.User

	// Query to retrieve the user by email
	row := config.DB.QueryRow("SELECT id, full_name, email, role, password FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("invalid email")
		}
		return user, err
	}

	// Validate the password
	isValid := helpers.CheckPasswordHashFunc(password, user.Password)
	if !isValid {
		return user, errors.New("invalid password")
	}

	return user, nil
}

func CheckUserAuthenticationMobile(email string, password string, branchID uint) (models.User, error) {
	var user models.User

	// Query to retrieve the user by email and branch_id
	row := config.DB.QueryRow(`
		SELECT id, full_name, email, role, password, branch_id 
		FROM users 
		WHERE email = $1 AND branch_id = $2`, email, branchID)

	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Password, &user.BranchId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("invalid email or branch ID")
		}
		return user, err
	}

	// Validate the password
	isValid := helpers.CheckPasswordHashFunc(password, user.Password)
	if !isValid {
		return user, errors.New("invalid password")
	}

	return user, nil
}
