package event

import (
	"context"
	"errors"

	shared "repo-backend/internal/content"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	List(context.Context, shared.ListOptions) ([]models.Event, int64, error)
	Find(context.Context, string, bool) (*models.Event, error)
	Create(context.Context, *models.Event) error
	Update(context.Context, *models.Event) error
	Delete(context.Context, string) error
}
type postgresRepository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &postgresRepository{db: db} }
func (r *postgresRepository) List(ctx context.Context, options shared.ListOptions) ([]models.Event, int64, error) {
	var items []models.Event
	query := r.db.WithContext(ctx).Model(&models.Event{}).Preload("AuthorProfile")
	if options.Published {
		query = query.Where("publication_status = ?", "PUBLISHED")
	}
	if options.Category != "" {
		query = query.Where("category = ?", options.Category)
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
	err := query.Order(shared.Order(options.Sort, "event_date")).
		Offset(options.Offset).Limit(options.Limit).Find(&items).Error
	return items, total, err
}
func (r *postgresRepository) Find(ctx context.Context, id string, published bool) (*models.Event, error) {
	var item models.Event
	query := r.db.WithContext(ctx).Preload("AuthorProfile")
	if published {
		query = query.Where("publication_status = ?", "PUBLISHED")
	}
	if err := query.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}
func (r *postgresRepository) Create(ctx context.Context, item *models.Event) error {
	return r.db.WithContext(ctx).Create(item).Error
}
func (r *postgresRepository) Update(ctx context.Context, item *models.Event) error {
	return r.db.WithContext(ctx).Omit("AuthorProfile").Save(item).Error
}
func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Event{}, "id = ?", id).Error
}
