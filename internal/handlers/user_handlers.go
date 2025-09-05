package handlers

import "gorm.io/gorm"

// LoginHandler manage user login
type LoginHandler struct {
	db *gorm.DB
}

// RegisterHandler manage user registration
type RegisterHandler struct {
	db *gorm.DB
}
