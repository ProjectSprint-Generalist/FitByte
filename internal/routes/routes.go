package routes

import (
	"fitbyte/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, userHandler *handlers.UserHandler, registerHandler *handlers.RegisterHandler, loginHandler *handlers.LoginHandler) {
	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Login & register routes
		v1.POST("/login", loginHandler.Login)
		v1.POST("/register", registerHandler.Register)

		// Health check routes
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.Health)
			health.GET("/ready", healthHandler.Ready)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("/", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
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
