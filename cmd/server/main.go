package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           app.Router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("failed to start server: ", err)
	}
}
