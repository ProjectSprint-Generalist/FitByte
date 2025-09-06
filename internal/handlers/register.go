package handlers

import (
	"errors"
	"fitbyte/internal/middleware"
	"fitbyte/internal/models"
	"fitbyte/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// NewRegisterHandler initializes RegisterHandler with the given DB
func NewRegisterHandler(db *gorm.DB) *RegisterHandler {
	return &RegisterHandler{
		db: db,
	}
}

// Register handle user registration requests
func (h *RegisterHandler) Register(context *gin.Context) {

	// Bind JSON to Register
	var inputUser models.InputUser
	if err := context.ShouldBindJSON(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please check your email and password format",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate email and password input
	if err := utils.Validate(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(inputUser.Password)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Failed to process password",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Create a new user
	user := &models.User{
		Email:    inputUser.Email,
		Password: hashedPassword,
	}

	if err := h.db.Create(user).Error; err != nil {
		// Check duplicate email
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			response := models.ErrorResponse{
				Success: false,
				Error:   "Email already registered",
				Code:    http.StatusConflict,
			}
			context.JSON(http.StatusConflict, response)
			return
		}
		response := models.ErrorResponse{
			Success: false,
			Error:   "An unexpected error occurred while creating the user",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Generate JWT Token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Failed to generate token",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"email": user.Email,
		"token": token,
	},
	)
}
