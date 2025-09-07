package routes

import (
	"fitbyte/internal/handlers"
	"fitbyte/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, userHandler *handlers.UserHandler, registerHandler *handlers.RegisterHandler, loginHandler *handlers.LoginHandler, activityHandler *handlers.ActivityHandler) {
	// API version 1
	v1 := router.Group("/v1")
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

		// User profile routes (auth required)
		userAuth := v1.Group("/user")
		userAuth.Use(middleware.IsAuthorized())
		{
			userAuth.GET("/", userHandler.GetUser)      
			userAuth.PATCH("/", userHandler.UpdateUser) 
		}

		activity := v1.Group("/activity")
		activity.Use(middleware.IsAuthorized())
		{
			activity.GET("/", activityHandler.GetActivities)
			activity.GET("/:id", activityHandler.GetActivity)
			activity.POST("/", activityHandler.CreateActivity)
			activity.DELETE("/:id", activityHandler.DeleteActivity)
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
