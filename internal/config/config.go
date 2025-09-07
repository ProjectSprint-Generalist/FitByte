package config

import (
	"fmt"
	"log"
	"os"

	"fitbyte/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds configuration values for the application
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
	DB          *gorm.DB
	MinIO       MinIOConfig
}

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func Load() *Config {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		MinIO: MinIOConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretAccessKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:          getEnv("MINIO_USE_SSL", "false") == "true",
			BucketName:      getEnv("MINIO_BUCKET_NAME", "fitbyte-files"),
		},
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
	err = db.AutoMigrate(&models.User{}, &models.FileUpload{})
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

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "fitbyte"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
	}
}

func (db *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db.Host, db.User, db.Password, db.Name, db.Port,
	)
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
