package models

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Data Transfer Object (DTO)
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// Validator
func (input *RegisterInput) Validate() error {

	// Email Validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		return errors.New("email format is invalid")
	}

	// Password Validation
	if len(input.Password) < 8 || len(input.Password) > 32 {
		return errors.New("password length must be 8–32 characters")
	}

	// Password Check Using Regex
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(input.Password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(input.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(input.Password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(input.Password)

	if !hasNumber || !hasUpper || !hasLower || !hasSpecial {
		return errors.New("password must contain at least one number, uppercase letter, lowercase letter, and special character")
	}
	return nil
}

// Hash Password
func (input *RegisterInput) HashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
