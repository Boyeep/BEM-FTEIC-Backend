package blog

import (
	"context"
	"time"

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

func (s *Service) List(ctx context.Context, query ListQuery, published bool) (*models.Page[models.Blog], error) {
	items, total, err := s.repo.List(ctx, shared.Options(query, published))
	if err != nil {
		return nil, shared.ServiceError(err, "blog")
	}
	return shared.Page(items, query.Page, query.PageSize, total), nil
}
func (s *Service) Get(ctx context.Context, id string, published bool) (*models.Blog, error) {
	item, err := s.repo.Find(ctx, id, published)
	if err != nil {
		return nil, shared.ServiceError(err, "blog")
	}
	return item, nil
}
func (s *Service) Create(ctx context.Context, input CreateDTO, userID string) (*models.Blog, error) {
	publishedAt := time.Now().UTC()
	if input.PublishedAt != "" {
		publishedAt, _ = time.Parse(time.RFC3339, input.PublishedAt)
	}
	creator := models.UUID(userID)
	item := &models.Blog{
		Title: input.Title, Excerpt: input.Excerpt, Author: input.Author,
		Category: input.Category, CoverImage: input.CoverImage, Content: input.Content,
		Status: input.Status, PublishedAt: publishedAt, CreatedBy: &creator,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, shared.ServiceError(err, "blog")
	}
	return item, nil
}
func (s *Service) Update(ctx context.Context, id string, input UpdateDTO) (*models.Blog, error) {
	item, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	publishedAt := time.Now().UTC()
	if input.PublishedAt != "" {
		publishedAt, _ = time.Parse(time.RFC3339, input.PublishedAt)
	}
	oldCover := item.CoverImage
	item.Title, item.Excerpt, item.Author = input.Title, input.Excerpt, input.Author
	item.Category, item.CoverImage, item.Content = input.Category, input.CoverImage, input.Content
	item.Status, item.PublishedAt = input.Status, publishedAt
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, shared.ServiceError(err, "blog")
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
		return shared.ServiceError(err, "blog")
	}
	if item.CoverImage != "" && s.media != nil {
		_ = s.media.Delete(ctx, item.CoverImage)
	}
	return nil
}
