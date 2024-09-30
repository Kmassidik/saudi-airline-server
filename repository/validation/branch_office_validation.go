package validation

import (
	"errors"
	"regexp"
)

// ValidateBranchOffices validates the input for the BranchOffice model
func ValidateBranchOffices(input map[string]interface{}) error {
	// Validate name
	if name, ok := input["name"].(string); !ok || name == "" {
		return errors.New("branch office name cannot be empty")
	}

	// Validate address
	if address, ok := input["address"].(string); !ok || address == "" {
		return errors.New("address cannot be empty")
	}

	// Validate total_counter
	var totalCounter float64
	if val, ok := input["total_counter"].(float64); !ok || val <= 0 {
		return errors.New("total counter must be greater than 0")
	} else {
		totalCounter = val
	}

	// Additional validations
	// 1. Validate that the name does not exceed a certain length
	if len(input["name"].(string)) > 255 {
		return errors.New("branch office name cannot exceed 255 characters")
	}

	// 2. Validate that the address does not exceed a certain length
	if len(input["address"].(string)) > 500 {
		return errors.New("address cannot exceed 500 characters")
	}

	// 3. Validate name format (e.g., only letters, numbers, and spaces)
	if !regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(input["name"].(string)) {
		return errors.New("branch office name can only contain letters, numbers, and spaces")
	}

	// 4. Ensure total_counter is not excessively large
	if totalCounter > 100 { // Adjust the limit as necessary
		return errors.New("total counter cannot exceed 100")
	}

	return nil
}
