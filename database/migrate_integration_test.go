//go:build integration

package database

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMigratePostgreSQL(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE SCHEMA IF NOT EXISTS auth;
		CREATE TABLE IF NOT EXISTS auth.users (id UUID PRIMARY KEY);`).Error; err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("migration must be idempotent: %v", err)
	}
	var count int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema='public' AND table_name='profiles' AND column_name='role'`).Scan(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected profiles.role column, got %d", count)
	}
}
