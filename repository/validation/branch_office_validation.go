package validation

import (
	"api-server/models"
	"errors"
)

// ValidateUser validates the input for the User model
func ValidateBranchOffices(user *models.BranchOffice) error {
	if user.Name == "" {
		return errors.New("name branch office cannot be empty")
	}

	if user.Address == "" {
		return errors.New("address cannot be empty")
	}

	if user.TotalCounter == 0 {
		return errors.New("total counter must be more then > 0")
	}

	return nil
}
