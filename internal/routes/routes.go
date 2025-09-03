package routes

import (
	"fitbyte/internal/handlers"
	"fitbyte/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, userHandler *handlers.UserHandler) {
	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Health check routes
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.Health)
			health.GET("/ready", healthHandler.Ready)
		}

		// User profile routes (auth required)
		userAuth := v1.Group("/user")
		userAuth.Use(middleware.DummyAuth()) // Use dummy auth middleware
		// userAuth.Use(middleware.Auth())
		{
			userAuth.GET("/", userHandler.GetUser)      // GET /v1/user
			userAuth.PATCH("/", userHandler.UpdateUser) // PATCH /v1/user
		}
	}

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to FitByte API",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})
}
