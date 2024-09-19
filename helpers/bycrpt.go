package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashingPasswordFunc hashes a plain text password
func HashingPasswordFunc(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword), nil
}

// CheckPasswordHashFunc compares a hashed password with a plain text one
func CheckPasswordHashFunc(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
