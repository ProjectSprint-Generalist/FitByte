package routes

import (
	"fitbyte/internal/handlers"
	"fitbyte/internal/middleware"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, activityHandler *handlers.ActivityHandler, jwtService *services.JWTService) {
	v1 := router.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.Health)
			health.GET("/ready", healthHandler.Ready)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			profile := protected.Group("/profile")
			{
				profile.GET("/", authHandler.Profile)
			}

			users := protected.Group("/users")
			{
				users.GET("/", userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.POST("/", userHandler.CreateUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			activities := protected.Group("/activities")
			{
				activities.GET("/", activityHandler.GetUserActivities)
				activities.GET("/all", activityHandler.GetActivities)
				activities.GET("/types", activityHandler.GetActivityTypes)
				activities.GET("/:id", activityHandler.GetActivity)
				activities.GET("/activity/:activityId", activityHandler.GetActivityByActivityID)
				activities.POST("/", activityHandler.CreateActivity)
				activities.PUT("/:id", activityHandler.UpdateActivity)
				activities.DELETE("/:id", activityHandler.DeleteActivity)
			}
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to FitByte API",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})
}
