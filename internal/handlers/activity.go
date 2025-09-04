package handlers

import (
	"net/http"
	"time"

	"fitbyte/internal/database"
	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
)

// ActivityHandler handles activity-related endpoints
type ActivityHandler struct{}

// NewActivityHandler creates a new activity handler
func NewActivityHandler() *ActivityHandler {
	return &ActivityHandler{}
}

// CreateActivity creates a new activity
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate activity type
	if !models.IsValidActivityType(string(req.ActivityType)) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity type. Must be one of: Walking, Yoga, Stretching, Cycling, Swimming, Dancing, Hiking, Running, HIIT, JumpRope",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate duration
	if req.DurationInMinutes < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Duration must be at least 1 minute",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Create activity
	activity := models.Activity{
		ActivityType:      req.ActivityType,
		DoneAt:            req.DoneAt,
		DurationInMinutes: req.DurationInMinutes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Calculate calories burned
	activity.CaloriesBurned = activity.CalculateCalories()

	// Save to database
	db := database.GetDB()
	if err := db.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to create activity",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Return response
	response := activity.ToResponse()
	c.JSON(http.StatusCreated, response)
}
