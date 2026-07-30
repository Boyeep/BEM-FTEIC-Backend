package service

import (
	"context"
	"errors"
	"time"

	"repo-backend/internal/dto"
	"repo-backend/internal/media"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"
	"repo-backend/pkg/apperr"
)

type BlogService struct {
	repo  repository.BlogRepository
	media *media.Service
}
type EventService struct {
	repo  repository.EventRepository
	media *media.Service
}
type GalleryService struct {
	repo  repository.GalleryRepository
	media *media.Service
}

func NewBlog(repo repository.BlogRepository, mediaService *media.Service) *BlogService {
	return &BlogService{repo: repo, media: mediaService}
}
func NewEvent(repo repository.EventRepository, mediaService *media.Service) *EventService {
	return &EventService{repo: repo, media: mediaService}
}
func NewGallery(repo repository.GalleryRepository, mediaService *media.Service) *GalleryService {
	return &GalleryService{repo: repo, media: mediaService}
}

func serviceError(err error, resource string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperr.NotFound(resource)
	}
	return apperr.Internal(err)
}

func (s *BlogService) List(ctx context.Context, query dto.ListQuery, published bool) (*models.Page[models.Blog], error) {
	v, total, err := s.repo.List(ctx, listOptions(query, published))
	if err != nil {
		return nil, serviceError(err, "blog")
	}
	return newPage(v, query.Page, query.PageSize, total), nil
}
func (s *BlogService) Get(ctx context.Context, id string) (*models.Blog, error) {
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, serviceError(err, "blog")
	}
	return v, nil
}
func (s *BlogService) GetPublic(ctx context.Context, id string) (*models.Blog, error) {
	v, err := s.repo.FindPublished(ctx, id)
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
	creatorID := models.UUID(userID)
	v := &models.Blog{Title: in.Title, Excerpt: in.Excerpt, Author: in.Author, Category: in.Category, CoverImage: in.CoverImage, Content: in.Content, Status: in.Status, PublishedAt: publishedAt, CreatedBy: &creatorID}
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
	oldCover := v.CoverImage
	v.Title, v.Excerpt, v.Author, v.Category, v.CoverImage, v.Content, v.Status, v.PublishedAt = in.Title, in.Excerpt, in.Author, in.Category, in.CoverImage, in.Content, in.Status, publishedAt
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, serviceError(err, "blog")
	}
	cleanupReplacedMedia(ctx, s.media, oldCover, v.CoverImage)
	return v, nil
}
func (s *BlogService) Delete(ctx context.Context, id string) error {
	item, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return serviceError(err, "blog")
	}
	cleanupMedia(ctx, s.media, item.CoverImage)
	return nil
}

func (s *EventService) List(ctx context.Context, query dto.ListQuery, published bool) (*models.Page[models.Event], error) {
	v, total, err := s.repo.List(ctx, listOptions(query, published))
	if err != nil {
		return nil, serviceError(err, "event")
	}
	return newPage(v, query.Page, query.PageSize, total), nil
}
func (s *EventService) Get(ctx context.Context, id string) (*models.Event, error) {
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, serviceError(err, "event")
	}
	return v, nil
}
func (s *EventService) GetPublic(ctx context.Context, id string) (*models.Event, error) {
	v, err := s.repo.FindPublished(ctx, id)
	if err != nil {
		return nil, serviceError(err, "event")
	}
	return v, nil
}
func (s *EventService) Create(ctx context.Context, in dto.CreateEvent, userID string) (*models.Event, error) {
	creatorID := models.UUID(userID)
	v := &models.Event{Title: in.Title, Description: in.Description, Author: in.Author, Category: in.Category, CoverImage: in.CoverImage, EventDate: models.Date(in.EventDate), Status: in.Status, PublicationStatus: in.PublicationStatus, CreatedBy: &creatorID}
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
	oldCover := v.CoverImage
	v.Title, v.Description, v.Author, v.Category, v.CoverImage, v.EventDate, v.Status, v.PublicationStatus = in.Title, in.Description, in.Author, in.Category, in.CoverImage, models.Date(in.EventDate), in.Status, in.PublicationStatus
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, serviceError(err, "event")
	}
	cleanupReplacedMedia(ctx, s.media, oldCover, v.CoverImage)
	return v, nil
}
func (s *EventService) Delete(ctx context.Context, id string) error {
	item, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return serviceError(err, "event")
	}
	cleanupMedia(ctx, s.media, item.CoverImage)
	return nil
}

func (s *GalleryService) List(ctx context.Context, query dto.ListQuery) (*models.Page[models.Gallery], error) {
	v, total, err := s.repo.List(ctx, listOptions(query, false))
	if err != nil {
		return nil, serviceError(err, "gallery item")
	}
	return newPage(v, query.Page, query.PageSize, total), nil
}
func (s *GalleryService) Get(ctx context.Context, id string) (*models.Gallery, error) {
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, serviceError(err, "gallery item")
	}
	return v, nil
}
func (s *GalleryService) Create(ctx context.Context, in dto.CreateGallery, userID string) (*models.Gallery, error) {
	if in.Category == "" {
		in.Category = "all"
	}
	creatorID := models.UUID(userID)
	v := &models.Gallery{Title: in.Title, Link: in.Link, ImageURL: in.ImageURL, Category: in.Category, TakenAt: models.Date(in.TakenAt), CreatedBy: &creatorID}
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
	if in.Category == "" {
		in.Category = v.Category
	}
	oldImage := v.ImageURL
	v.Title, v.Link, v.ImageURL, v.Category, v.TakenAt = in.Title, in.Link, in.ImageURL, in.Category, models.Date(in.TakenAt)
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, serviceError(err, "gallery item")
	}
	cleanupReplacedMedia(ctx, s.media, oldImage, v.ImageURL)
	return v, nil
}
func (s *GalleryService) Delete(ctx context.Context, id string) error {
	item, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return serviceError(err, "gallery item")
	}
	cleanupMedia(ctx, s.media, item.ImageURL)
	return nil
}

func cleanupReplacedMedia(ctx context.Context, mediaService *media.Service, oldURL, newURL string) {
	if oldURL != newURL {
		cleanupMedia(ctx, mediaService, oldURL)
	}
}

func cleanupMedia(ctx context.Context, mediaService *media.Service, rawURL string) {
	if mediaService != nil && rawURL != "" {
		_ = mediaService.Delete(ctx, rawURL)
	}
}

func offset(page, pageSize int) int {
	return (page - 1) * pageSize
}

func listOptions(query dto.ListQuery, published bool) repository.ListOptions {
	return repository.ListOptions{
		Offset:    offset(query.Page, query.PageSize),
		Limit:     query.PageSize,
		Category:  query.Category,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
		Sort:      query.Sort,
		Published: published,
	}
}

func newPage[T any](items []T, page, pageSize int, total int64) *models.Page[T] {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	return &models.Page[T]{
		Items: items,
		Pagination: models.Pagination{
			Page:            page,
			PageSize:        pageSize,
			TotalItems:      total,
			TotalPages:      totalPages,
			HasNextPage:     page < totalPages,
			HasPreviousPage: page > 1,
		},
	}
}
