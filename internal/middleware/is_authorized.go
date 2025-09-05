package middleware

import (
	"fitbyte/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)
// Authorization
func IsAuthorized() gin.HandlerFunc {
	return func(context *gin.Context) {

		// Check for authorization header
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			response := models.ErrorResponse{
				Success: false,
				Error:   "Authorization header required",
				Code:    http.StatusUnauthorized,
			}
			context.JSON(http.StatusUnauthorized, response)
			context.Abort()
			return
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
		context.Set("userID", claims.ID)
		context.Set("email", claims.Email)

		context.Next()
	}
}
