package handlers

import (
	"net/http"
	"strconv"
	"time"

	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
)

// ActivityHandler handles activity-related endpoints
type ActivityHandler struct {

	// In a real application, you would inject a service or repository here
	// activityService services.ActivityService
}

// NewActivityHandler creates a new activity handler
func NewActivityHandler() *ActivityHandler {
	return &ActivityHandler{}
}

// GetUsers returns a list of users
func (h *ActivityHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	actId1 := 1
	actId2 := 2
	acttype1 := "Activity 1"
	acttype2 := "Activity 2"
	durationInM1 := 20
	durationInM2 := 30
	calories1 := 100
	calories2 := 200
	users := []models.Activity{
		{
			BaseEntity: models.BaseEntity{
				ID:        actId1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ActivityType:      acttype1,
			DurationInMinutes: durationInM1,
			CaloriesBurned:    calories1,
		},
		{
			BaseEntity: models.BaseEntity{
				ID:        actId2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ActivityType:      acttype2,
			DurationInMinutes: durationInM2,
			CaloriesBurned:    calories2,
		},
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int64(len(users)),
			TotalPages: 1,
		},
	})
}

// GetUser returns a specific user by ID
func (h *ActivityHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	actId1 := id
	acttype1 := "Activity 1"
	durationInM1 := 20
	calories1 := 100
	activity := models.Activity{
		BaseEntity: models.BaseEntity{
			ID:        int(actId1),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ActivityType:      acttype1,
		DurationInMinutes: durationInM1,
		CaloriesBurned:    calories1,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    activity,
	})
}
