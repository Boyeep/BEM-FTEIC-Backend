package main

import (
	"repo-backend/internal/database"
	"repo-backend/internal/models"
	"repo-backend/internal/routes"
)

func main() {
	// Connect to database
	database.ConnectDB()

	// Migrate the schema
	database.DB.AutoMigrate(&models.Blog{})

	// Setup router
	r := routes.SetupRouter()

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	// Run server
	r.Run(":8080")
}
