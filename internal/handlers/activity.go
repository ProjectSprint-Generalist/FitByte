package handlers

import (
	"net/http"
	"strconv"
	"time"

	"fitbyte/internal/database"
	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateActivity creates a new activity

type ActivityHandler struct {
	db *gorm.DB
}

func NewActivityHandler(db *gorm.DB) *ActivityHandler {
	return &ActivityHandler{db: db}
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req models.CreateActivityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if !models.IsValidActivityType(string(req.ActivityType)) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Invalid activity type. Must be one of: Walking, Yoga, Stretching, Cycling, Swimming, Dancing, Hiking, Running, HIIT, JumpRope",
			})

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
				})
			}
			// Return response
			response := activity.ToResponse()
			c.JSON(http.StatusCreated, response)
		}
	}
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	query := h.db.Model(&models.Activity{})

	if activityType := c.Query("activityType"); activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	if doneAtFromStr := c.Query("doneAtFrom"); doneAtFromStr != "" {
		if doneAtFrom, err := time.Parse(time.RFC3339, doneAtFromStr); err == nil {
			query = query.Where("done_at >= ?", doneAtFrom)
		}
	}

	if doneAtToStr := c.Query("doneAtTo"); doneAtToStr != "" {
		if doneAtTo, err := time.Parse(time.RFC3339, doneAtToStr); err == nil {
			query = query.Where("done_at <= ?", doneAtTo)
		}
	}

	if caloriesMinStr := c.Query("caloriesBurnedMin"); caloriesMinStr != "" {
		if caloriesMin, err := strconv.Atoi(caloriesMinStr); err == nil {
			query = query.Where("calories_burned >= ?", caloriesMin)
		}
	}

	if caloriesMaxStr := c.Query("caloriesBurnedMax"); caloriesMaxStr != "" {
		if caloriesMax, err := strconv.Atoi(caloriesMaxStr); err == nil {
			query = query.Where("calories_burned <= ?", caloriesMax)
		}
	}

	var activities []models.Activity
	var total int64

	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Message: "Activities retrieved successfully",
		Data:    activities,
		Pagination: models.Pagination{
			Page:       (offset / limit) + 1,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// GetActivity returns a specific activity by ID
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var activity models.Activity
	if err := h.db.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Activity not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Activity retrieved successfully",
		Data:    activity,
	})
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // Convert :id parameter to integer
	if err != nil {
		// Return a 400 Bad Request response if the ID is invalid
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil { // If the request body is invalid (invalid JSON)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var activity models.Activity
	if err := h.db.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Activity not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	if req.ActivityType == nil && req.DurationInMinutes == nil && req.CaloriesBurned == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Apply updates
	if req.ActivityType != nil {
		activity.ActivityType = *req.ActivityType
	}
	if req.DurationInMinutes != nil {
		activity.DurationInMinutes = *req.DurationInMinutes
	}
	if req.CaloriesBurned != nil {
		activity.CaloriesBurned = *req.CaloriesBurned
	}

	if err := h.db.Save(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Activity updated successfully",
		Data:    activity,
	})
}
