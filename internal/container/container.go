package container

import (
	"repo-backend/config"
	"repo-backend/internal/blog"
	"repo-backend/internal/event"
	"repo-backend/internal/gallery"
	"repo-backend/internal/handlers"
	"repo-backend/internal/media"
	"repo-backend/internal/repository"
	"repo-backend/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	Accounts  *handlers.Account
	Analytics *handlers.Analytics
	Blog      *blog.Handler
	Event     *event.Handler
	Gallery   *gallery.Handler
	Media     *handlers.Media
}

func NewContainer(db *gorm.DB) *Container {
	mediaService := media.NewService(
		media.NewLocalStorage("uploads"),
		config.Optional("PUBLIC_API_URL", ""),
	)
	return &Container{
		Accounts:  &handlers.Account{Service: service.NewAccount(repository.NewAccountRepository(db))},
		Analytics: &handlers.Analytics{Service: service.NewAnalytics(repository.NewAnalyticsRepository(db))},
		Blog:      blog.NewHandler(blog.NewService(blog.NewRepository(db), mediaService)),
		Event:     event.NewHandler(event.NewService(event.NewRepository(db), mediaService)),
		Gallery:   gallery.NewHandler(gallery.NewService(gallery.NewRepository(db), mediaService)),
		Media:     &handlers.Media{Service: mediaService},
	}
}
