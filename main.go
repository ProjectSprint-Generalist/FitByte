package main

import (
	"log"
	"os"

	"fitbyte/internal/config"
	"fitbyte/internal/database"
	"fitbyte/internal/handlers"
	"fitbyte/internal/middleware"
	"fitbyte/internal/routes"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	userService := services.NewUserService()
	jwtService := services.NewJWTService(cfg)

	healthHandler := handlers.NewHealthHandler()
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	routes.SetupRoutes(router, healthHandler, userHandler, authHandler, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
