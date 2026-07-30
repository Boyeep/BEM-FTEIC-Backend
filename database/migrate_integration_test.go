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
	if err := db.Exec(`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='anon') THEN CREATE ROLE anon; END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='authenticated') THEN CREATE ROLE authenticated; END IF;
		END $$;
		CREATE SCHEMA IF NOT EXISTS auth;
		CREATE TABLE IF NOT EXISTS auth.users (
			id UUID PRIMARY KEY,
			email TEXT,
			raw_user_meta_data JSONB NOT NULL DEFAULT '{}'::JSONB
		);
		CREATE OR REPLACE FUNCTION auth.uid() RETURNS UUID
		LANGUAGE sql STABLE AS 'SELECT NULL::UUID';`).Error; err != nil {
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
	var policyCount int64
	if err := db.Raw(`SELECT COUNT(*) FROM pg_policies
		WHERE schemaname='public' AND tablename IN ('profiles','blogs','events','galeri','signup_whitelist')`).Scan(&policyCount).Error; err != nil {
		t.Fatal(err)
	}
	if policyCount < 8 {
		t.Fatalf("expected security policies, got %d", policyCount)
	}
	var supportTableCount int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema='public' AND table_name IN ('signup_whitelist','site_visitors')`).Scan(&supportTableCount).Error; err != nil {
		t.Fatal(err)
	}
	if supportTableCount != 2 {
		t.Fatalf("expected support tables, got %d", supportTableCount)
	}
	var publicationColumnCount int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema='public' AND table_name='events'
		  AND column_name='publication_status'`).Scan(&publicationColumnCount).Error; err != nil {
		t.Fatal(err)
	}
	if publicationColumnCount != 1 {
		t.Fatalf("expected events.publication_status column, got %d", publicationColumnCount)
	}
	var eventPolicy string
	if err := db.Raw(`SELECT qual FROM pg_policies
		WHERE schemaname='public' AND tablename='events'
		  AND policyname='events_public_read'`).Scan(&eventPolicy).Error; err != nil {
		t.Fatal(err)
	}
	if eventPolicy == "" {
		t.Fatal("expected published-only event policy")
	}
}
