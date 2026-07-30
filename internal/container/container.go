package container

import (
	"repo-backend/config"
	"repo-backend/internal/handlers"
	"repo-backend/internal/media"
	"repo-backend/internal/repository"
	"repo-backend/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	Accounts  *handlers.Account
	Analytics *handlers.Analytics
	Content   *handlers.Content
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
		Content: &handlers.Content{
			Blogs:   service.NewBlog(repository.NewContentRepository(db), mediaService),
			Events:  service.NewEvent(repository.NewEventRepository(db), mediaService),
			Gallery: service.NewGallery(repository.NewGalleryRepository(db), mediaService),
		},
		Media: &handlers.Media{Service: mediaService},
	}
}
