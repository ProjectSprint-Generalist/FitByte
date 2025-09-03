package config

import (
	"fmt"
	"log"
	"os"

	"fitbyte/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds all configuration for our application
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
	DB          *gorm.DB
}

// Load reads configuration from environment variables and initializes database
func Load() *Config {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
	}

	// Initialize database
	cfg.initDatabase()

	return cfg
}

// initDatabase initializes the database connection
func (c *Config) initDatabase() {
	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(c.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default user if not exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		defaultUser := &models.User{
			ID:    1,
			Email: "test@example.com",
			Name:  stringPtr("Test User"),
		}
		db.Create(defaultUser)
		fmt.Println("Default user created with ID: 1")
	}

	c.DB = db
	fmt.Println("Database connected successfully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
