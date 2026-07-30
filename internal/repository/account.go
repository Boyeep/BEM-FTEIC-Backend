package repository

import (
	"context"
	"errors"
	"time"

	"repo-backend/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AccountRepository interface {
	EnsureProfile(context.Context, string, string, string) error
	FindProfile(context.Context, string) (*models.Profile, error)
	UpdateProfile(context.Context, string, string, string) error
	ListPublicProfiles(context.Context, []string) ([]models.PublicProfile, error)
	ListWhitelist(context.Context) ([]models.WhitelistEntry, error)
	CreateWhitelist(context.Context, *models.WhitelistEntry) error
	DeleteWhitelist(context.Context, string) error
}

type accountRepository struct{ db *gorm.DB }

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) EnsureProfile(ctx context.Context, id, email, username string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO profiles (id,email,username,role)
		VALUES (?,?,?,'member') ON CONFLICT (id) DO NOTHING`, id, email, username).Error
}

func (r *accountRepository) FindProfile(ctx context.Context, id string) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.WithContext(ctx).Table("profiles").
		Select("id,email,username,avatar_url,role").
		Where("id = ?", id).Take(&profile).Error
	return &profile, translateError(err)
}

func (r *accountRepository) UpdateProfile(ctx context.Context, id, username, avatarURL string) error {
	result := r.db.WithContext(ctx).Table("profiles").Where("id = ?", id).
		Updates(map[string]any{
			"username": username, "avatar_url": avatarURL, "updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *accountRepository) ListPublicProfiles(ctx context.Context, ids []string) ([]models.PublicProfile, error) {
	var profiles []models.PublicProfile
	err := r.db.WithContext(ctx).Table("profiles").
		Select("id,username,avatar_url").
		Where("id IN ?", ids).Find(&profiles).Error
	return profiles, err
}

func (r *accountRepository) ListWhitelist(ctx context.Context) ([]models.WhitelistEntry, error) {
	var rows []models.WhitelistEntry
	err := r.db.WithContext(ctx).Table("signup_whitelist").
		Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *accountRepository) CreateWhitelist(ctx context.Context, entry *models.WhitelistEntry) error {
	err := r.db.WithContext(ctx).Table("signup_whitelist").Create(entry).Error
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return ErrConflict
	}
	return err
}

func (r *accountRepository) DeleteWhitelist(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Exec("DELETE FROM signup_whitelist WHERE id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
