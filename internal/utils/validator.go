package utils

import (
	"errors"
	"fitbyte/internal/models"
	"regexp"
)

// Validator
func Validate(input *models.InputUser) error {

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
