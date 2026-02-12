package database

import (
	"fmt"
	"log"
	"os"

	"repo-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Connection string based on your credentials
	// Note: We use 'bem_backend' as the database name. You need to create this DB first!
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.Admin{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Connected to PostgreSQL database")
}
