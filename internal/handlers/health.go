package handlers

import (
	"net/http"

	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
)


type HealthHandler struct{}


func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}


func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is healthy",
		Data: gin.H{
			"status":    "ok",
			"service":   "fitbyte-api",
			"version":   "1.0.0",
			"timestamp": gin.H{},
		},
	})
}


func (h *HealthHandler) Ready(c *gin.Context) {

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is ready",
		Data: gin.H{
			"status": "ready",
		},
	})
}
