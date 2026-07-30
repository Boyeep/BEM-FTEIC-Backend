package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"repo-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "admin"),
		getEnv("DB_NAME", "bem_backend"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
		getEnv("DB_TIMEZONE", "Asia/Jakarta"),
	)

	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Database connection attempt %d/10 failed: %v", attempt, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to database after 10 attempts:", err)
	}

	err = DB.AutoMigrate(&models.Admin{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB.AutoMigrate(&models.Blog{})
	DB.AutoMigrate(&models.Gallery{})
	DB.AutoMigrate(&models.Event{})

	log.Println("Connected to PostgreSQL database")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
