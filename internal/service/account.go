package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"repo-backend/internal/models"
	"repo-backend/internal/repository"
	"repo-backend/pkg/apperr"
)

type AccountService struct{ repo repository.AccountRepository }

func NewAccount(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Me(ctx context.Context, id, email string) (*models.Profile, error) {
	username := strings.Split(email, "@")[0]
	if err := s.repo.EnsureProfile(ctx, id, email, username); err != nil {
		return nil, apperr.Internal(err)
	}
	profile, err := s.repo.FindProfile(ctx, id)
	if err != nil {
		return nil, serviceError(err, "profile")
	}
	return profile, nil
}

func (s *AccountService) UpdateProfile(ctx context.Context, id, username, avatarURL string) (*models.Profile, error) {
	if err := s.repo.UpdateProfile(ctx, id, strings.TrimSpace(username), avatarURL); err != nil {
		return nil, serviceError(err, "profile")
	}
	profile, err := s.repo.FindProfile(ctx, id)
	if err != nil {
		return nil, serviceError(err, "profile")
	}
	return profile, nil
}

func (s *AccountService) PublicProfiles(ctx context.Context, rawIDs string) ([]models.PublicProfile, error) {
	seen := make(map[string]struct{})
	ids := make([]string, 0)
	for _, rawID := range strings.Split(rawIDs, ",") {
		id := strings.TrimSpace(rawID)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
		if len(ids) == 50 {
			break
		}
	}
	if len(ids) == 0 {
		return []models.PublicProfile{}, nil
	}
	profiles, err := s.repo.ListPublicProfiles(ctx, ids)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return profiles, nil
}

func (s *AccountService) ListWhitelist(ctx context.Context) ([]models.WhitelistEntry, error) {
	rows, err := s.repo.ListWhitelist(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return rows, nil
}

func (s *AccountService) AddWhitelist(ctx context.Context, email, creator string) (*models.WhitelistEntry, error) {
	creatorID := models.UUID(creator)
	entry := &models.WhitelistEntry{
		Email:     strings.ToLower(strings.TrimSpace(email)),
		CreatedBy: &creatorID,
	}
	if err := s.repo.CreateWhitelist(ctx, entry); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, apperr.New("CONFLICT", "email is already whitelisted", http.StatusConflict)
		}
		return nil, apperr.Internal(err)
	}
	return entry, nil
}

func (s *AccountService) DeleteWhitelist(ctx context.Context, id string) error {
	if err := s.repo.DeleteWhitelist(ctx, id); err != nil {
		return serviceError(err, "whitelist entry")
	}
	return nil
}
