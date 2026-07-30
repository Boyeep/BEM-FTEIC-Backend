package service

import (
	"context"
	"errors"
	"time"

	"repo-backend/internal/dto"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"
	"repo-backend/pkg/apperr"

	"gorm.io/gorm"
)

type BlogService struct{ repo repository.BlogRepository }
type EventService struct{ repo repository.EventRepository }
type GalleryService struct{ repo repository.GalleryRepository }

func NewBlog(repo repository.BlogRepository) *BlogService    { return &BlogService{repo: repo} }
func NewEvent(repo repository.EventRepository) *EventService { return &EventService{repo: repo} }
func NewGallery(repo repository.GalleryRepository) *GalleryService {
	return &GalleryService{repo: repo}
}

func serviceError(err error, resource string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.NotFound(resource)
	}
	return apperr.Internal(err)
}

func (s *BlogService) List(ctx context.Context, published bool) ([]models.Blog, error) {
	v, err := s.repo.List(ctx, published)
	if err != nil {
		return nil, serviceError(err, "blog")
	}
	return v, nil
}
func (s *BlogService) Get(ctx context.Context, id string) (*models.Blog, error) {
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, serviceError(err, "blog")
	}
	return v, nil
}
func (s *BlogService) Create(ctx context.Context, in dto.CreateBlog, userID string) (*models.Blog, error) {
	publishedAt := time.Now().UTC()
	if in.PublishedAt != "" {
		publishedAt, _ = time.Parse(time.RFC3339, in.PublishedAt)
	}
	v := &models.Blog{Title: in.Title, Excerpt: in.Excerpt, Author: in.Author, Category: in.Category, CoverImage: in.CoverImage, Content: in.Content, Status: in.Status, PublishedAt: publishedAt, CreatedBy: &userID}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, serviceError(err, "blog")
	}
	return v, nil
}
func (s *BlogService) Update(ctx context.Context, id string, in dto.UpdateBlog) (*models.Blog, error) {
	v, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	publishedAt := time.Now().UTC()
	if in.PublishedAt != "" {
		publishedAt, _ = time.Parse(time.RFC3339, in.PublishedAt)
	}
	v.Title, v.Excerpt, v.Author, v.Category, v.CoverImage, v.Content, v.Status, v.PublishedAt = in.Title, in.Excerpt, in.Author, in.Category, in.CoverImage, in.Content, in.Status, publishedAt
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, serviceError(err, "blog")
	}
	return v, nil
}
func (s *BlogService) Delete(ctx context.Context, id string) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return serviceError(err, "blog")
	}
	return nil
}

func (s *EventService) List(ctx context.Context, category string, published bool) ([]models.Event, error) {
	v, err := s.repo.List(ctx, category, published)
	if err != nil {
		return nil, serviceError(err, "event")
	}
	return v, nil
}
func (s *EventService) Get(ctx context.Context, id string) (*models.Event, error) {
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, serviceError(err, "event")
	}
	return v, nil
}
func (s *EventService) Create(ctx context.Context, in dto.CreateEvent, userID string) (*models.Event, error) {
	v := &models.Event{Title: in.Title, Description: in.Description, Author: in.Author, Category: in.Category, CoverImage: in.CoverImage, EventDate: in.EventDate, Status: in.Status, CreatedBy: &userID}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, serviceError(err, "event")
	}
	return v, nil
}
func (s *EventService) Update(ctx context.Context, id string, in dto.UpdateEvent) (*models.Event, error) {
	v, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	v.Title, v.Description, v.Author, v.Category, v.CoverImage, v.EventDate, v.Status = in.Title, in.Description, in.Author, in.Category, in.CoverImage, in.EventDate, in.Status
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, serviceError(err, "event")
	}
	return v, nil
}
func (s *EventService) Delete(ctx context.Context, id string) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return serviceError(err, "event")
	}
	return nil
}

func (s *GalleryService) List(ctx context.Context) ([]models.Gallery, error) {
	v, err := s.repo.List(ctx)
	if err != nil {
		return nil, serviceError(err, "gallery item")
	}
	return v, nil
}
func (s *GalleryService) Get(ctx context.Context, id string) (*models.Gallery, error) {
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, serviceError(err, "gallery item")
	}
	return v, nil
}
func (s *GalleryService) Create(ctx context.Context, in dto.CreateGallery, userID string) (*models.Gallery, error) {
	v := &models.Gallery{Title: in.Title, Link: in.Link, ImageURL: in.ImageURL, TakenAt: in.TakenAt, CreatedBy: &userID}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, serviceError(err, "gallery item")
	}
	return v, nil
}
func (s *GalleryService) Update(ctx context.Context, id string, in dto.UpdateGallery) (*models.Gallery, error) {
	v, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	v.Title, v.Link, v.ImageURL, v.TakenAt = in.Title, in.Link, in.ImageURL, in.TakenAt
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, serviceError(err, "gallery item")
	}
	return v, nil
}
func (s *GalleryService) Delete(ctx context.Context, id string) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return serviceError(err, "gallery item")
	}
	return nil
}
