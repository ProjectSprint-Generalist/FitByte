package handlers

import (
	"net/http"
	"strconv"

	"fitbyte/internal/models"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService *services.ActivityService
}

func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	activityType := c.Query("type")

	var activities []models.Activity
	var total int64
	var err error

	if activityType != "" {
		activities, total, err = h.activityService.GetActivitiesByType(activityType, page, limit)
	} else {
		activities, total, err = h.activityService.GetAllActivities(page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	var activityResponses []models.ActivityResponse
	for _, activity := range activities {
		activityResponses = append(activityResponses, h.activityService.ToActivityResponse(&activity))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Activities retrieved successfully",
		Data: gin.H{
			"activities": activityResponses,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

func (h *ActivityHandler) GetUserActivities(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	activityType := c.Query("type")

	var activities []models.Activity
	var total int64
	var err error

	if activityType != "" {
		activities, total, err = h.activityService.GetUserActivitiesByType(userID.(uint), activityType, page, limit)
	} else {
		activities, total, err = h.activityService.GetUserActivities(userID.(uint), page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	var activityResponses []models.ActivityResponse
	for _, activity := range activities {
		activityResponses = append(activityResponses, h.activityService.ToActivityResponse(&activity))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User activities retrieved successfully",
		Data: gin.H{
			"activities": activityResponses,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

func (h *ActivityHandler) GetActivity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	activity, err := h.activityService.GetActivityByID(uint(id))
	if err != nil {
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
		Data:    h.activityService.ToActivityResponse(activity),
	})
}

func (h *ActivityHandler) GetActivityByActivityID(c *gin.Context) {
	activityID := c.Param("activityId")

	activity, err := h.activityService.GetActivityByActivityID(activityID)
	if err != nil {
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
		Data:    h.activityService.ToActivityResponse(activity),
	})
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
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
			Error:   "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	activity, err := h.activityService.CreateActivity(userID.(uint), &req)
	if err != nil {
		if err.Error() == "validation error: invalid activity type" || 
		   err.Error() == "invalid activity type" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Activity created successfully",
		Data:    h.activityService.ToActivityResponse(activity),
	})
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	activity, err := h.activityService.UpdateActivity(uint(id), &req)
	if err != nil {
		if err.Error() == "validation error: invalid activity type" || 
		   err.Error() == "invalid activity type" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		
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
		Data:    h.activityService.ToActivityResponse(activity),
	})
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid activity ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.activityService.DeleteActivity(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Activity deleted successfully",
		Data:    nil,
	})
}

func (h *ActivityHandler) GetActivityTypes(c *gin.Context) {
	activityTypes := h.activityService.GetActivityTypes()
	
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Activity types retrieved successfully",
		Data: gin.H{
			"activityTypes": activityTypes,
		},
	})
}