package validation

import (
	"api-server/models"
	"errors"
	"regexp"
)

// ValidateUser validates the input for the User model
func ValidateUser(user *models.User) error {
	if user.FullName == "" {
		return errors.New("full name cannot be empty")
	}

	if user.Email == "" {
		return errors.New("email cannot be empty")
	}

	if user.Password == "" {
		return errors.New("password cannot be empty")
	}

	if len(user.FullName) < 3 {
		return errors.New("full name must be at least 3 characters")
	}

	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if user.Role != "administrator" && user.Role != "admin" && user.Role != "supervisor" && user.Role != "officer" {
		return errors.New("role must be either 'administrator' 'admin', 'supervisor', or 'officer'")
	}

	return nil
}

func CheckLoginUserInput(email string, password string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	if password == "" {
		return errors.New("password cannot be empty")
	}

	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}
