package validation

import (
	"errors"
)

// ValidateBranchCounter validates the input data for the BranchCounter model
func ValidateBranchCounter(input map[string]interface{}) error {
	// Validate "counter_location" must be a non-empty string
	counterLocation, ok := input["counter_location"].(string)
	if !ok || counterLocation == "" {
		return errors.New("counter_location must be a non-empty string")
	}

	// Validate "user_id" must be a number (convertible to uint
	userIDStr, ok := input["branch_id"].(float64)

	if !ok {
		return errors.New("user_id must be a non-empty string")
	}

	// Use branchIDFloat or convert it to uint if needed
	_ = uint(userIDStr)

	// Validate "branch_id" must be a number (convertible to uint)
	branchIDFloat, ok := input["branch_id"].(float64)
	if !ok {
		return errors.New("branch_id must be a valid number")
	}

	// Use branchIDFloat or convert it to uint if needed
	_ = uint(branchIDFloat)

	// Everything is valid
	return nil
}
