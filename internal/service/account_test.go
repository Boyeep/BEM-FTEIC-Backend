package service

import (
	"context"
	"errors"
	"testing"

	"repo-backend/internal/models"
	"repo-backend/internal/repository"
)

type accountRepositoryStub struct {
	createErr error
	created   *models.WhitelistEntry
	deletedID string
	deleteErr error
}

func (s *accountRepositoryStub) EnsureProfile(context.Context, string, string, string) error {
	return nil
}
func (s *accountRepositoryStub) FindProfile(context.Context, string) (*models.Profile, error) {
	return &models.Profile{}, nil
}
func (s *accountRepositoryStub) UpdateProfile(context.Context, string, string, string) error {
	return nil
}
func (s *accountRepositoryStub) ListPublicProfiles(context.Context, []string) ([]models.PublicProfile, error) {
	return nil, nil
}
func (s *accountRepositoryStub) ListWhitelist(context.Context) ([]models.WhitelistEntry, error) {
	return nil, nil
}
func (s *accountRepositoryStub) CreateWhitelist(_ context.Context, entry *models.WhitelistEntry) error {
	s.created = entry
	return s.createErr
}
func (s *accountRepositoryStub) DeleteWhitelist(_ context.Context, id string) error {
	s.deletedID = id
	return s.deleteErr
}

func TestAccountAddWhitelistNormalizesEmail(t *testing.T) {
	repo := &accountRepositoryStub{}
	entry, err := NewAccount(repo).AddWhitelist(context.Background(), " Admin@Example.COM ", "creator")
	if err != nil {
		t.Fatal(err)
	}
	if entry.Email != "admin@example.com" || repo.created == nil {
		t.Fatalf("unexpected entry: %#v", entry)
	}
}

func TestAccountAddWhitelistMapsConflict(t *testing.T) {
	repo := &accountRepositoryStub{createErr: repository.ErrConflict}
	_, err := NewAccount(repo).AddWhitelist(context.Background(), "admin@example.com", "creator")
	if err == nil || !errors.Is(repo.createErr, repository.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestAccountDeleteWhitelistMapsNotFound(t *testing.T) {
	repo := &accountRepositoryStub{deleteErr: repository.ErrNotFound}
	if err := NewAccount(repo).DeleteWhitelist(context.Background(), "missing"); err == nil {
		t.Fatal("expected not found error")
	}
}
