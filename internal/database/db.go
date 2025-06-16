package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	maxRetries := 30
	baseDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err == nil {
			fmt.Println("Database connected successfully")
			return
		}

		delay := time.Duration(i+1) * baseDelay
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, delay)
		time.Sleep(delay)
	}

	log.Fatal("Failed to connect to database after", maxRetries, "attempts:", err)
}

func GetDB() *gorm.DB {
	return DB
}
