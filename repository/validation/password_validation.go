package validation

import (
	"errors"
	"regexp"
)

// ValidatePassword checks the password against common policies.
func ValidatePassword(password string) error {
	// Password must be at least 6 characters long
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// At least one uppercase letter
	upperCaseRegex := `[A-Z]`
	if !regexp.MustCompile(upperCaseRegex).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// At least one lowercase letter
	lowerCaseRegex := `[a-z]`
	if !regexp.MustCompile(lowerCaseRegex).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// At least one digit
	digitRegex := `[0-9]`
	if !regexp.MustCompile(digitRegex).MatchString(password) {
		return errors.New("password must contain at least one digit")
	}

	// At least one special character
	specialCharRegex := `[!@#\$%^&*()_+{}\[\]:;"'<>,.?/\\|]`
	if !regexp.MustCompile(specialCharRegex).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
