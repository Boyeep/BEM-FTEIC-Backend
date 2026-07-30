package event

import (
	"context"

	shared "repo-backend/internal/content"
	"repo-backend/internal/media"
	"repo-backend/internal/models"
)

type Service struct {
	repo  Repository
	media *media.Service
}

func NewService(repo Repository, mediaService *media.Service) *Service {
	return &Service{repo: repo, media: mediaService}
}
func (s *Service) List(ctx context.Context, query ListQuery, published bool) (*models.Page[models.Event], error) {
	items, total, err := s.repo.List(ctx, shared.Options(query, published))
	if err != nil {
		return nil, shared.ServiceError(err, "event")
	}
	return shared.Page(items, query.Page, query.PageSize, total), nil
}
func (s *Service) Get(ctx context.Context, id string, published bool) (*models.Event, error) {
	item, err := s.repo.Find(ctx, id, published)
	if err != nil {
		return nil, shared.ServiceError(err, "event")
	}
	return item, nil
}
func (s *Service) Create(ctx context.Context, input CreateDTO, userID string) (*models.Event, error) {
	creator := models.UUID(userID)
	item := &models.Event{
		Title: input.Title, Description: input.Description, Author: input.Author,
		Category: input.Category, CoverImage: input.CoverImage,
		EventDate: models.Date(input.EventDate), Status: input.Status,
		PublicationStatus: input.PublicationStatus, CreatedBy: &creator,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, shared.ServiceError(err, "event")
	}
	return item, nil
}
func (s *Service) Update(ctx context.Context, id string, input UpdateDTO) (*models.Event, error) {
	item, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	oldCover := item.CoverImage
	item.Title, item.Description, item.Author = input.Title, input.Description, input.Author
	item.Category, item.CoverImage = input.Category, input.CoverImage
	item.EventDate, item.Status = models.Date(input.EventDate), input.Status
	item.PublicationStatus = input.PublicationStatus
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, shared.ServiceError(err, "event")
	}
	if oldCover != item.CoverImage && oldCover != "" && s.media != nil {
		_ = s.media.Delete(ctx, oldCover)
	}
	return item, nil
}
func (s *Service) Delete(ctx context.Context, id string) error {
	item, err := s.Get(ctx, id, false)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return shared.ServiceError(err, "event")
	}
	if item.CoverImage != "" && s.media != nil {
		_ = s.media.Delete(ctx, item.CoverImage)
	}
	return nil
}
