package repository

import (
	"context"

	"repo-backend/internal/models"

	"gorm.io/gorm"
)

type BlogRepository interface {
	List(context.Context, bool) ([]models.Blog, error)
	Find(context.Context, string) (*models.Blog, error)
	Create(context.Context, *models.Blog) error
	Update(context.Context, *models.Blog) error
	Delete(context.Context, string) error
}

type EventRepository interface {
	List(context.Context, string, bool) ([]models.Event, error)
	Find(context.Context, string) (*models.Event, error)
	Create(context.Context, *models.Event) error
	Update(context.Context, *models.Event) error
	Delete(context.Context, string) error
}

type GalleryRepository interface {
	List(context.Context) ([]models.Gallery, error)
	Find(context.Context, string) (*models.Gallery, error)
	Create(context.Context, *models.Gallery) error
	Update(context.Context, *models.Gallery) error
	Delete(context.Context, string) error
}

type contentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *contentRepository {
	return &contentRepository{db: db}
}

func (r *contentRepository) List(ctx context.Context, published bool) ([]models.Blog, error) {
	var items []models.Blog
	query := r.db.WithContext(ctx).Order("published_at DESC")
	if published {
		query = query.Where("status = ?", "PUBLISHED")
	}
	return items, query.Find(&items).Error
}

func (r *contentRepository) Find(ctx context.Context, id string) (*models.Blog, error) {
	var item models.Blog
	return &item, r.db.WithContext(ctx).First(&item, "id = ?", id).Error
}

func (r *contentRepository) Create(ctx context.Context, item *models.Blog) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *contentRepository) Update(ctx context.Context, item *models.Blog) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *contentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Blog{}, "id = ?", id).Error
}

type eventRepository struct{ db *gorm.DB }

func NewEventRepository(db *gorm.DB) EventRepository { return &eventRepository{db: db} }
func (r *eventRepository) List(ctx context.Context, category string, published bool) ([]models.Event, error) {
	var items []models.Event
	query := r.db.WithContext(ctx).Order("event_date DESC")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if published {
		query = query.Where("status = ?", "PUBLISHED")
	}
	return items, query.Find(&items).Error
}
func (r *eventRepository) Find(ctx context.Context, id string) (*models.Event, error) {
	var item models.Event
	return &item, r.db.WithContext(ctx).First(&item, "id = ?", id).Error
}
func (r *eventRepository) Create(ctx context.Context, item *models.Event) error {
	return r.db.WithContext(ctx).Create(item).Error
}
func (r *eventRepository) Update(ctx context.Context, item *models.Event) error {
	return r.db.WithContext(ctx).Save(item).Error
}
func (r *eventRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Event{}, "id = ?", id).Error
}

type galleryRepository struct{ db *gorm.DB }

func NewGalleryRepository(db *gorm.DB) GalleryRepository { return &galleryRepository{db: db} }
func (r *galleryRepository) List(ctx context.Context) ([]models.Gallery, error) {
	var items []models.Gallery
	return items, r.db.WithContext(ctx).Order("taken_at DESC").Find(&items).Error
}
func (r *galleryRepository) Find(ctx context.Context, id string) (*models.Gallery, error) {
	var item models.Gallery
	return &item, r.db.WithContext(ctx).First(&item, "id = ?", id).Error
}
func (r *galleryRepository) Create(ctx context.Context, item *models.Gallery) error {
	return r.db.WithContext(ctx).Create(item).Error
}
func (r *galleryRepository) Update(ctx context.Context, item *models.Gallery) error {
	return r.db.WithContext(ctx).Save(item).Error
}
func (r *galleryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Gallery{}, "id = ?", id).Error
}
