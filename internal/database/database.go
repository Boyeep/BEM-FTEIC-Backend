package database

import (
	"log"

	"repo-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Connection string based on your credentials
	// Note: We use 'bem_backend' as the database name. You need to create this DB first!
	dsn := "host=localhost user=postgres password=admin dbname=bem_backend port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
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
