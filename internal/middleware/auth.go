package middleware

import (
	"fitbyte/internal/models"
	"log"
	"os"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Generate Token for Login and Register
func GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
        // log.Println("JWT_SECRET not set in environment")
        return "", fmt.Errorf("JWT_SECRET not set")
    }
	
	expTime := time.Now().Add(30*time.Minute)

	claims := &JWTClaim{
		ID : user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject: string(rune(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
        log.Println("Failed to sign JWT:", err)
        return "", err
    }
	// if secret == "" {
    //     return "", fmt.Errorf("JWT_SECRET not set")
    // }
	
	return signed, nil
}