package handlers

import (
	"net/http"
	"strconv"

	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActivityHandler struct {
	db *gorm.DB
}

func NewActivityHandler(db *gorm.DB) *ActivityHandler {
	return &ActivityHandler{db: db}
}

// GetActivities returns a list of activities
func (h *ActivityHandler) GetActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var activities []models.Activity
	var total int64

	h.db.Model(&models.Activity{}).Count(&total)

	if err := h.db.Limit(limit).Offset(offset).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Message: "Activities retrieved successfully",
		Data:    activities,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
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
