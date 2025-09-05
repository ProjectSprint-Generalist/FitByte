package handlers

import "gorm.io/gorm"

// RegisterHandler manage user registration
type RegisterHandler struct {
	db *gorm.DB
}
