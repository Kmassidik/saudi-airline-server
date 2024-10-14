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
			"SELECT id, full_name, email, role, likes, dislikes, image, branch_id FROM users WHERE role = ? ORDER BY id ASC LIMIT ? OFFSET ?",
			role, limit, offset,
		)
	} else {
		rows, err = config.DB.Query(
			"SELECT id, full_name, email, role, likes, dislikes, image, branch_id FROM users WHERE role != ? ORDER BY id ASC LIMIT ? OFFSET ?",
			"officer", limit, offset,
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
	row := config.DB.QueryRow("SELECT id, full_name, email, role, likes, dislikes, image, branch_id, password FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Likes, &user.Dislikes, &user.Image, &user.BranchId, &user.Password)

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

	// Begin a transaction to ensure atomicity
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// Defer rollback to ensure any error during execution rolls back the transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert the user into the database, including the image path
	_, err = tx.Exec(
		"INSERT INTO users (full_name, email, password, role, likes, dislikes, image, branch_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.FullName, user.Email, user.Password, user.Role, user.Likes, user.Dislikes, user.Image, user.BranchId,
	)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}

	// Update total officer count only if the user's role is "officer"
	if user.Role == "officer" {
		_, err = tx.Exec("UPDATE total_data SET total_officer = total_officer + 1 WHERE id = 1")
		if err != nil {
			log.Println("Error updating total_officer:", err)
			return err
		}
	}

	// Commit the transaction if no errors
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
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

	// Check if the branch exists before proceeding
	exists, err := GetBranchOfficesById(user.BranchId)
	if err != nil {
		return err
	}

	if exists == nil {
		return errors.New("branch_id does not exist")
	}

	// Retrieve the current user data (to compare roles)
	var currentRole string
	err = config.DB.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&currentRole)
	if err != nil {
		log.Println("Error retrieving current user role:", err)
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

	// Begin a transaction to ensure atomicity
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// Defer rollback to ensure any error during execution rolls back the transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare the SQL query and execute
	if user.Image != "" {
		query := "UPDATE users SET full_name = ?, email = ?, password = ?, role = ?, image = ?, branch_id = ? WHERE id = ?"
		_, err = tx.Exec(query, user.FullName, user.Email, user.Password, user.Role, user.Image, user.BranchId, id)
	} else {
		query := "UPDATE users SET full_name = ?, email = ?, password = ?, role = ?, branch_id = ? WHERE id = ?"
		_, err = tx.Exec(query, user.FullName, user.Email, user.Password, user.Role, user.BranchId, id)
	}

	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	// Compare old role with the new role and update total_officer accordingly
	if currentRole != "officer" && user.Role == "officer" {
		// If user is promoted to officer
		_, err = tx.Exec("UPDATE total_data SET total_officer = total_officer + 1 WHERE id = 1")
		if err != nil {
			log.Println("Error increasing total_officer:", err)
			return err
		}
	} else if currentRole == "officer" && user.Role != "officer" {
		// If user is demoted from officer
		_, err = tx.Exec("UPDATE total_data SET total_officer = total_officer - 1 WHERE id = 1")
		if err != nil {
			log.Println("Error decreasing total_officer:", err)
			return err
		}
	}

	// Commit the transaction if no errors
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil
}

// DeleteUser removes a user from the database and decrements the total officer count if the user is an officer
func DeleteUser(id uint) error {
	// Begin a transaction to ensure atomicity
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// Defer rollback to ensure any error during execution rolls back the transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the user's role before deletion
	var role string
	err = tx.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&role)
	if err != nil {
		log.Println("Error fetching user role:", err)
		return err
	}

	// Delete the user from the users table
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	// Decrement total officer count only if the user's role is "officer"
	if role == "officer" {
		_, err = tx.Exec("UPDATE total_data SET total_officer = total_officer - 1 WHERE id = 1")
		if err != nil {
			log.Println("Error updating total_officer:", err)
			return err
		}
	}

	// Commit the transaction if no errors
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil
}

func GetAllUsersByBranchOfiice(branchId uint) ([]models.UserByBranchOfiiceResponse, error) {
	var users []models.UserByBranchOfiiceResponse

	// Execute a SELECT query to fetch users by branch_id
	rows, err := config.DB.Query("SELECT id, full_name FROM users WHERE branch_id = ? AND role = 'officer' AND NOT EXISTS ( SELECT 1 FROM branch_counters WHERE branch_counters.user_id = users.id);", branchId)
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	// Scan the result rows into the users slice
	for rows.Next() {
		var user models.UserByBranchOfiiceResponse
		if err := rows.Scan(&user.ID, &user.FullName); err != nil {
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

// CheckUserAuthentication checks if the user exists with the given email and password.
func CheckUserAuthentication(email string, password string) (models.User, error) {
	var user models.User

	// Query to retrieve the user by email
	query := "SELECT id, full_name, email, role, password FROM users WHERE email = ?"
	row := config.DB.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("invalid email")
		}
		return user, err
	}

	// Validate the password
	if !helpers.CheckPasswordHashFunc(password, user.Password) {
		return user, errors.New("invalid password")
	}

	return user, nil
}

// CheckUserAuthenticationMobile checks if the user exists with the given email, password, and branch ID.
func CheckUserAuthenticationMobile(email string, password string, branchID uint) (models.User, error) {
	var user models.User

	// Query to retrieve the user by email, including branch ID
	query := `
		SELECT id, full_name, email, role, password, branch_id 
		FROM users 
		WHERE email = ?`
	err := config.DB.QueryRow(query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Password, &user.BranchId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("email not found")
		}
		return user, err
	}

	// Validate the branch ID
	if user.BranchId != branchID {
		return user, errors.New("invalid branch office")
	}

	// Validate the password
	if !helpers.CheckPasswordHashFunc(password, user.Password) {
		log.Printf("Invalid password for user: %s", email)
		return user, errors.New("invalid password")
	}

	return user, nil
}
