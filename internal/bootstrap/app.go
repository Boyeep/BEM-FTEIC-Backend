package bootstrap

import (
	"fmt"

	"repo-backend/config"
	"repo-backend/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func New() (*App, error) {
	db, err := config.OpenDatabase()
	if err != nil {
		return nil, fmt.Errorf("bootstrap database: %w", err)
	}
	router, err := routes.SetupRouter(db)
	if err != nil {
		return nil, fmt.Errorf("bootstrap router: %w", err)
	}

	return &App{DB: db, Router: router}, nil
}
