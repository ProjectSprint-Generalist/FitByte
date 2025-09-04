package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash Password
func HashPassword(inputPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// Verify Password
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
