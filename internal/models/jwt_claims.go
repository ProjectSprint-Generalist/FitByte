package models

import "github.com/golang-jwt/jwt/v5"

// JWT Payload
type JWTClaim struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
