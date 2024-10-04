package services

import (
	"api-server/models"
	"api-server/repository"
	"errors"
)

// GetAllUsers retrieves all users with pagination
func GetAllUsers(limit, offset int, role string) ([]models.UserAllResponse, error) {
	return repository.GetAllUsers(limit, offset, role)
}

// GetUserByID retrieves a user by ID
func GetUserByID(id uint) (*models.User, error) {
	return repository.GetUserByID(id)
}

// CreateUser creates a new user
func CreateUser(user *models.User) error {
	// Perform validation before insertion
	if err := repository.CreateUser(user); err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user
func UpdateUser(id uint, user *models.User) error {
	return repository.UpdateUser(id, user)
}

// DeleteUser deletes a user by ID
func DeleteUser(id uint) error {
	return repository.DeleteUser(id)
}

// GetUsersCount fetches the count of users based on the role from the repository.
func GetUsersCount(role string) (int, error) {
	// Call the repository function and get the count and error
	count, err := repository.GetUsersCount(role)
	if err != nil {
		// If there is an error, return 0 and the error
		return 0, err
	}

	// Return the user count and nil for error if successful
	return count, nil
}

func GetUsersByBranchID(brancdId uint) ([]models.UserByBranchOfiiceResponse, error) {
	return repository.GetAllUsersByBranchOfiice(brancdId)
}

// Authentication
func AuthenticationLoginUser(email string, password string) (models.User, error) {
	// Call the repository function to check authentication
	user, err := repository.CheckUserAuthentication(email, password)
	if err != nil {
		return user, err
	}

	// If user is not an admin, return an authorization error
	if user.Role != "administrator" && user.Role != "admin" && user.Role != "supervisor" {
		return user, errors.New("authorization failed: user is not authorized")
	}

	return user, nil
}

func AuthenticationLoginUserMobile(email string, password string, branch_id uint) (models.User, error) {
	// Call the repository function to check authentication
	user, err := repository.CheckUserAuthenticationMobile(email, password, branch_id)
	if err != nil {
		return user, err
	}

	// If user is not an admin, return an authorization error
	if user.Role != "admin" && user.Role != "supervisor" {
		return user, errors.New("authorization failed: user is not authorized")
	}

	return user, nil
}
