package database

import (
	"fitbyte/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to Database
func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to Database")
	}

	// Assign to global DB
	DB = db

	// Connection pool
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	log.Println("Connected to Database")
}

// Auto Migration
func Migrate() {
	DB.AutoMigrate(&models.User{})
	log.Println("Database Migration Completed")
}
