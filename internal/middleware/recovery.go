package middleware

import (
	"net/http"

	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
)


func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Internal server error: " + err,
				Code:    http.StatusInternalServerError,
			})
		}
		c.Abort()
	})
}
