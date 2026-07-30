package blog

import (
	"context"
	"errors"

	shared "repo-backend/internal/content"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	List(context.Context, shared.ListOptions) ([]models.Blog, int64, error)
	Find(context.Context, string, bool) (*models.Blog, error)
	Create(context.Context, *models.Blog) error
	Update(context.Context, *models.Blog) error
	Delete(context.Context, string) error
}

type postgresRepository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &postgresRepository{db: db} }

func (r *postgresRepository) List(ctx context.Context, options shared.ListOptions) ([]models.Blog, int64, error) {
	var items []models.Blog
	query := r.db.WithContext(ctx).Model(&models.Blog{}).Preload("AuthorProfile")
	if options.Published {
		query = query.Where("status = ?", "PUBLISHED")
	}
	if options.Category != "" {
		query = query.Where("category = ?", options.Category)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order(shared.Order(options.Sort, "published_at")).
		Offset(options.Offset).Limit(options.Limit).Find(&items).Error
	return items, total, err
}

func (r *postgresRepository) Find(ctx context.Context, id string, published bool) (*models.Blog, error) {
	var item models.Blog
	query := r.db.WithContext(ctx).Preload("AuthorProfile")
	if published {
		query = query.Where("status = ?", "PUBLISHED")
	}
	if err := query.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepository) Create(ctx context.Context, item *models.Blog) error {
	return r.db.WithContext(ctx).Create(item).Error
}
func (r *postgresRepository) Update(ctx context.Context, item *models.Blog) error {
	return r.db.WithContext(ctx).Save(item).Error
}
func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Blog{}, "id = ?", id).Error
}
