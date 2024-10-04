package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Load JWT secret from environment variable
var jwtKey = []byte(os.Getenv("JWT_TOKEN"))

// Custom Claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new token for a valid user
func GenerateJWT(username string) (string, error) {
	// Set token expiration time, e.g., 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT parses and validates the JWT token string
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token and validate the signature
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
