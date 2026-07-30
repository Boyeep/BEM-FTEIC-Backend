package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	Track(context.Context, string, string, string) error
	Count(context.Context) (int64, error)
}

type analyticsRepository struct{ db *gorm.DB }

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) Track(ctx context.Context, id, path, userAgent string) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO site_visitors (id,last_seen_at,last_path,user_agent)
		VALUES (?,?,?,?)
		ON CONFLICT (id) DO UPDATE SET
			last_seen_at=EXCLUDED.last_seen_at,
			last_path=EXCLUDED.last_path,
			user_agent=EXCLUDED.user_agent`,
		id, time.Now().UTC(), path, userAgent,
	).Error
}

func (r *analyticsRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("site_visitors").Count(&count).Error
	return count, err
}
