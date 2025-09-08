package handlers

import (
	"net/http"

	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// GetUser returns current user profile (GET /v1/user)
func (h *UserHandler) GetUser(c *gin.Context) {
	// Get user ID from auth middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Find user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "User not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Database error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert to response format
	userResponse := models.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		ImageURI:   user.ImageURI,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data:    userResponse,
	})
}

// UpdateUser updates user profile (PATCH /v1/user)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Get user ID from auth middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Parse request body
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation error: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find existing user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "User not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Database error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Preference != nil {
		updates["preference"] = *req.Preference
	}
	if req.WeightUnit != nil {
		updates["weight_unit"] = *req.WeightUnit
	}
	if req.HeightUnit != nil {
		updates["height_unit"] = *req.HeightUnit
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	}
	if req.ImageURI != nil {
		updates["image_uri"] = *req.ImageURI
	}

	// Perform update
	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to update user",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Reload user data
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to reload user data",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert to response format
	userResponse := models.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		ImageURI:   user.ImageURI,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User profile updated successfully",
		Data:    userResponse,
	})
}
