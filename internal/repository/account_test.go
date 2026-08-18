package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"repo-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func accountRepositoryDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db, mock
}

func TestCreateWhitelistUsesDatabaseGeneratedFields(t *testing.T) {
	db, mock := accountRepositoryDatabase(t)
	createdAt := time.Now().UTC()
	creator := models.UUID("11111111-1111-4111-8111-111111111111")
	entry := &models.WhitelistEntry{
		Email:     "website.fteic@gmail.com",
		CreatedBy: &creator,
	}
	mock.ExpectQuery(`INSERT INTO signup_whitelist \(email, created_by\).*RETURNING id, email, created_by, created_at`).
		WithArgs(entry.Email, string(creator)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "created_by", "created_at"}).AddRow(
			"22222222-2222-4222-8222-222222222222",
			entry.Email,
			string(creator),
			createdAt,
		))

	if err := NewAccountRepository(db).CreateWhitelist(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	if entry.ID == "" || entry.CreatedAt.IsZero() {
		t.Fatalf("expected database-generated fields, got %#v", entry)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateWhitelistMapsDuplicateEmail(t *testing.T) {
	db, mock := accountRepositoryDatabase(t)
	creator := models.UUID("11111111-1111-4111-8111-111111111111")
	entry := &models.WhitelistEntry{
		Email:     "website.fteic@gmail.com",
		CreatedBy: &creator,
	}
	mock.ExpectQuery(`INSERT INTO signup_whitelist \(email, created_by\)`).
		WithArgs(entry.Email, string(creator)).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	err := NewAccountRepository(db).CreateWhitelist(context.Background(), entry)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
