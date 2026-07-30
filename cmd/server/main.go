package main

import (
	"log"
	"os"

	"repo-backend/internal/database"
	"repo-backend/internal/routes"
)

func main() {
	// Connect to database
	database.ConnectDB()

	// Setup router
	r := routes.SetupRouter()

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
