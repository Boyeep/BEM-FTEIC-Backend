package main

import (
	"log"
	"os"

	"repo-backend/config"
	"repo-backend/database"
	"repo-backend/internal/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "--migrate" {
		if err := database.Migrate(app.DB); err != nil {
			log.Fatal("database migration failed: ", err)
		}
		log.Println("database migrations completed")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "--schema" {
		if err := database.PrintSchema(app.DB, os.Stdout); err != nil {
			log.Fatal("database schema inspection failed: ", err)
		}
		return
	}

	port := config.Optional("PORT", "8080")
	if err := app.Router.Run(":" + port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
