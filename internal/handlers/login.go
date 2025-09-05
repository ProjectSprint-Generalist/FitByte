package handlers

import (
	"fitbyte/internal/middleware"
	"fitbyte/internal/models"
	"fitbyte/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewLoginHandler initializes LoginHandler with the given DB
func NewLoginHandler(db *gorm.DB) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

// Login handle user login requests
func (h *LoginHandler) Login(context *gin.Context) {

	var inputUser models.InputUser

	// Bind JSON
	if err := context.ShouldBindJSON(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please provide a valid email and password",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	//  Check existing user
	var user models.User
	if err := h.db.Where("email = ?", inputUser.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := models.ErrorResponse{
				Success: false,
				Error:   "User does not exist",
				Code:    http.StatusNotFound,
			}
			context.JSON(http.StatusNotFound, response)
			return
		}
		response := models.ErrorResponse{
			Success: false,
			Error:   "Server error",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Verify Password
	if err := utils.VerifyPassword(inputUser.Password, user.Password); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid password",
			Code:    http.StatusUnauthorized,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}

	// Generate Token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Failed to generate token",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Login successfully
	context.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User login successfully",
		Data: gin.H{
			"email": user.Email,
			"token": token,
		},
	})
}
