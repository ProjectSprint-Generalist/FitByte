package handlers

import (
	"net/http"
	"time"

	"fitbyte/internal/database"
	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		db: database.DB,
	}
}

// Health returns the health status of the API
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is healthy",
		Data: gin.H{
			"status":    "ok",
			"service":   "fitbyte-api",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
		},
	})
}

// Ready returns the readiness status of the API
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check database connection
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
				Success: false,
				Error:   "Database connection error",
				Code:    http.StatusServiceUnavailable,
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
				Success: false,
				Error:   "Database ping failed",
				Code:    http.StatusServiceUnavailable,
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is ready",
		Data: gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
		},
	})
}
