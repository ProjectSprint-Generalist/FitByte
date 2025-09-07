package handlers

import (
	"net/http"
	"strconv"
	"time"

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
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if !models.IsValidActivityType(string(req.ActivityType)) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity type. Must be one of: Walking, Yoga, Stretching, Cycling, Swimming, Dancing, Hiking, Running, HIIT, JumpRope",
			Code:    http.StatusBadRequest,
		})
		return
	}

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
		UserID:            userID.(uint), // Set user ID from JWT token
		ActivityType:      req.ActivityType,
		DoneAt:            req.DoneAt,
		DurationInMinutes: req.DurationInMinutes,
	}

	// Calculate calories burned
	activity.CaloriesBurned = activity.CalculateCalories()

	// Save to database
	if err := h.db.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to create activity",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Activity created successfully",
		Data:    activity.ToResponse(),
	})
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

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

	// Filter by user ID for security
	query := h.db.Model(&models.Activity{}).Where("user_id = ?", userID)

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

	// Convert to response format
	var activityResponses []models.ActivityResponse
	for _, activity := range activities {
		activityResponses = append(activityResponses, activity.ToResponse())
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Message: "Activities retrieved successfully",
		Data:    activityResponses,
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
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	id, err := strconv.Atoi(c.Param("activityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var activity models.Activity
	// Filter by both ID and user_id for security
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&activity).Error; err != nil {
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
		Data:    activity.ToResponse(),
	})
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	id, err := strconv.Atoi(c.Param("activityId")) // Convert :activityId parameter to integer
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
			Error:   "Request binding error: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var activity models.Activity
	// Filter by both ID and user_id for security
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Activity not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	if req.ActivityType == nil && req.DurationInMinutes == nil && req.DoneAt == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "At least one field must be provided for update",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Apply updates
	if req.ActivityType != nil {
		if !models.IsValidActivityType(*req.ActivityType) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Invalid activity type. Must be one of: Walking, Yoga, Stretching, Cycling, Swimming, Dancing, Hiking, Running, HIIT, JumpRope",
				Code:    http.StatusBadRequest,
			})
			return
		}
		activity.ActivityType = models.ActivityType(*req.ActivityType)
	}
	if req.DurationInMinutes != nil {
		if *req.DurationInMinutes < 1 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Duration must be at least 1 minute",
				Code:    http.StatusBadRequest,
			})
			return
		}
		activity.DurationInMinutes = *req.DurationInMinutes
		// Recalculate calories when duration changes
		activity.CaloriesBurned = activity.CalculateCalories()
	}
	if req.DoneAt != nil {
		activity.DoneAt = *req.DoneAt
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
		Data:    activity.ToResponse(),
	})
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	id, err := strconv.Atoi(c.Param("activityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var activity models.Activity

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Activity not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	if err := h.db.Delete(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to delete activity",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Activity deleted successfully",
	})
}
