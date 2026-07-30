package gallery

import (
	"context"
	"errors"

	shared "repo-backend/internal/content"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	List(context.Context, shared.ListOptions) ([]models.Gallery, int64, error)
	Find(context.Context, string) (*models.Gallery, error)
	Create(context.Context, *models.Gallery) error
	Update(context.Context, *models.Gallery) error
	Delete(context.Context, string) error
}
type postgresRepository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &postgresRepository{db: db} }
func (r *postgresRepository) List(ctx context.Context, options shared.ListOptions) ([]models.Gallery, int64, error) {
	var items []models.Gallery
	query := r.db.WithContext(ctx).Model(&models.Gallery{}).Preload("AuthorProfile")
	if options.Category != "" && options.Category != "all" {
		query = query.Where("category = ?", options.Category)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order(shared.Order(options.Sort, "taken_at")).
		Offset(options.Offset).Limit(options.Limit).Find(&items).Error
	return items, total, err
}
func (r *postgresRepository) Find(ctx context.Context, id string) (*models.Gallery, error) {
	var item models.Gallery
	if err := r.db.WithContext(ctx).Preload("AuthorProfile").First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}
func (r *postgresRepository) Create(ctx context.Context, item *models.Gallery) error {
	return r.db.WithContext(ctx).Create(item).Error
}
func (r *postgresRepository) Update(ctx context.Context, item *models.Gallery) error {
	return r.db.WithContext(ctx).Omit("AuthorProfile").Save(item).Error
}
func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Gallery{}, "id = ?", id).Error
}
