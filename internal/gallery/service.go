package gallery

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
func (s *Service) List(ctx context.Context, query ListQuery) (*models.Page[models.Gallery], error) {
	items, total, err := s.repo.List(ctx, shared.Options(query, false))
	if err != nil {
		return nil, shared.ServiceError(err, "gallery item")
	}
	return shared.Page(items, query.Page, query.PageSize, total), nil
}
func (s *Service) Get(ctx context.Context, id string) (*models.Gallery, error) {
	item, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, shared.ServiceError(err, "gallery item")
	}
	return item, nil
}
func (s *Service) Create(ctx context.Context, input CreateDTO, userID string) (*models.Gallery, error) {
	if input.Category == "" {
		input.Category = "all"
	}
	creator := models.UUID(userID)
	item := &models.Gallery{
		Title: input.Title, Link: input.Link, ImageURL: input.ImageURL,
		Category: input.Category, TakenAt: models.Date(input.TakenAt), CreatedBy: &creator,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, shared.ServiceError(err, "gallery item")
	}
	return item, nil
}
func (s *Service) Update(ctx context.Context, id string, input UpdateDTO) (*models.Gallery, error) {
	item, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Category == "" {
		input.Category = item.Category
	}
	oldImage := item.ImageURL
	item.Title, item.Link, item.ImageURL = input.Title, input.Link, input.ImageURL
	item.Category, item.TakenAt = input.Category, models.Date(input.TakenAt)
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, shared.ServiceError(err, "gallery item")
	}
	if oldImage != item.ImageURL && oldImage != "" && s.media != nil {
		_ = s.media.Delete(ctx, oldImage)
	}
	return item, nil
}
func (s *Service) Delete(ctx context.Context, id string) error {
	item, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return shared.ServiceError(err, "gallery item")
	}
	if item.ImageURL != "" && s.media != nil {
		_ = s.media.Delete(ctx, item.ImageURL)
	}
	return nil
}
