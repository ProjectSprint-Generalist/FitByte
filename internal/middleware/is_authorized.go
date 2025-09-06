package middleware

import (
	"fitbyte/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authorization
func IsAuthorized() gin.HandlerFunc {
	return func(context *gin.Context) {

		// Check for authorization header
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			response := models.ErrorResponse{
				Success: false,
				Error:   "Authorization header required",
				Code:    http.StatusUnauthorized,
			}
			context.JSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Parse and validate JWT
		claims, err := ParseToken(tokenString)
		if err != nil {
			response := models.ErrorResponse{
				Success: false,
				Error:   "Invalid or expired token",
				Code:    http.StatusUnauthorized,
			}
			context.JSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}

		// Store user info
		context.Set("user_id", claims.ID)
		context.Set("email", claims.Email)

		context.Next()
	}
}
