package services

import (
	"api-server/models"
	"api-server/repository"
)

// GetAllUsers retrieves all users with pagination
func GetAllUsers(limit, offset int) ([]models.UserResponse, error) {
	return repository.GetAllUsers(limit, offset)
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
