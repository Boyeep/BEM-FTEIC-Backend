package repository

import (
	"context"
	"errors"

	"repo-backend/internal/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("repository: not found")
var ErrConflict = errors.New("repository: conflict")

type ListOptions struct {
	Offset    int
	Limit     int
	Category  string
	StartDate string
	EndDate   string
	Sort      string
	Published bool
}

type BlogRepository interface {
	List(context.Context, ListOptions) ([]models.Blog, int64, error)
	Find(context.Context, string) (*models.Blog, error)
	FindPublished(context.Context, string) (*models.Blog, error)
	Create(context.Context, *models.Blog) error
	Update(context.Context, *models.Blog) error
	Delete(context.Context, string) error
}

type EventRepository interface {
	List(context.Context, ListOptions) ([]models.Event, int64, error)
	Find(context.Context, string) (*models.Event, error)
	FindPublished(context.Context, string) (*models.Event, error)
	Create(context.Context, *models.Event) error
	Update(context.Context, *models.Event) error
	Delete(context.Context, string) error
}

type GalleryRepository interface {
	List(context.Context, ListOptions) ([]models.Gallery, int64, error)
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

func (r *contentRepository) List(ctx context.Context, options ListOptions) ([]models.Blog, int64, error) {
	var items []models.Blog
	query := r.db.WithContext(ctx).Model(&models.Blog{})
	if options.Published {
		query = query.Where("status = ?", "PUBLISHED")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := listOrder(options.Sort, "published_at")
	err := query.Order(order).Offset(options.Offset).Limit(options.Limit).Find(&items).Error
	return items, total, err
}

func (r *contentRepository) Find(ctx context.Context, id string) (*models.Blog, error) {
	var item models.Blog
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	return &item, translateError(err)
}
func (r *contentRepository) FindPublished(ctx context.Context, id string) (*models.Blog, error) {
	var item models.Blog
	err := r.db.WithContext(ctx).Where("status = ?", "PUBLISHED").First(&item, "id = ?", id).Error
	return &item, translateError(err)
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
func (r *eventRepository) List(ctx context.Context, options ListOptions) ([]models.Event, int64, error) {
	var items []models.Event
	query := r.db.WithContext(ctx).Model(&models.Event{})
	if options.Category != "" {
		query = query.Where("category = ?", options.Category)
	}
	if options.Published {
		query = query.Where("publication_status = ?", "PUBLISHED")
	}
	if options.StartDate != "" {
		query = query.Where("event_date >= ?", options.StartDate)
	}
	if options.EndDate != "" {
		query = query.Where("event_date <= ?", options.EndDate)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order(listOrder(options.Sort, "event_date")).
		Offset(options.Offset).Limit(options.Limit).Find(&items).Error
	return items, total, err
}
func (r *eventRepository) Find(ctx context.Context, id string) (*models.Event, error) {
	var item models.Event
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	return &item, translateError(err)
}
func (r *eventRepository) FindPublished(ctx context.Context, id string) (*models.Event, error) {
	var item models.Event
	err := r.db.WithContext(ctx).Where("publication_status = ?", "PUBLISHED").First(&item, "id = ?", id).Error
	return &item, translateError(err)
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
func (r *galleryRepository) List(ctx context.Context, options ListOptions) ([]models.Gallery, int64, error) {
	var items []models.Gallery
	query := r.db.WithContext(ctx).Model(&models.Gallery{})
	if options.Category != "" && options.Category != "all" {
		query = query.Where("category = ?", options.Category)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order(listOrder(options.Sort, "taken_at")).
		Offset(options.Offset).Limit(options.Limit).Find(&items).Error
	return items, total, err
}
func (r *galleryRepository) Find(ctx context.Context, id string) (*models.Gallery, error) {
	var item models.Gallery
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	return &item, translateError(err)
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

func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

func listOrder(sort, dateColumn string) string {
	switch sort {
	case "oldest":
		return dateColumn + " ASC"
	case "title_asc":
		return "title ASC"
	case "title_desc":
		return "title DESC"
	default:
		return dateColumn + " DESC"
	}
}
