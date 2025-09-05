package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// DB is the database connection pool
	DB *gorm.DB
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	var err error

	// config obj can be filled with logger
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	// turn this on to enable connection pooling for connection adjustment
	// pgDb, err := DB.DB()
	// if err != nil {
	// 	log.Fatal("Failed to get database connection:", err)
	// 	return nil, err
	// }

	// // Connection pooling
	// pgDb.SetMaxIdleConns(10)
	// pgDb.SetMaxOpenConns(100)

	log.Println("Connected to database")

	return DB, nil
}
