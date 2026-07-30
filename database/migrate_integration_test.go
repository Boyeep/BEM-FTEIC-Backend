//go:build integration

package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"repo-backend/internal/blog"
	shared "repo-backend/internal/content"
	"repo-backend/internal/event"
	"repo-backend/internal/gallery"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"

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
	if err := db.Exec(`INSERT INTO auth.users (id,email,raw_user_meta_data)
		VALUES ('55555555-5555-4555-8555-555555555555', 'stevenprobot@gmail.com',
			'{"username":"Project Owner"}'::JSONB)
		ON CONFLICT (id) DO UPDATE
		SET email=EXCLUDED.email, raw_user_meta_data=EXCLUDED.raw_user_meta_data`).Error; err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("migration must be idempotent: %v", err)
	}
	var ownerRole string
	if err := db.Raw(
		"SELECT role FROM profiles WHERE email = ?",
		"stevenprobot@gmail.com",
	).Scan(&ownerRole).Error; err != nil {
		t.Fatal(err)
	}
	if ownerRole != "admin" {
		t.Fatalf("expected bootstrapped owner to be admin, got %q", ownerRole)
	}
	testVerticalRepositories(t, db)
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

func testVerticalRepositories(t *testing.T, db *gorm.DB) {
	t.Helper()
	ctx := context.Background()
	creator := models.UUID("11111111-1111-4111-8111-111111111111")
	blogID := models.UUID("22222222-2222-4222-8222-222222222222")
	draftBlogID := models.UUID("22222222-2222-4222-8222-222222222223")
	eventID := models.UUID("33333333-3333-4333-8333-333333333333")
	draftEventID := models.UUID("33333333-3333-4333-8333-333333333334")
	galleryID := models.UUID("44444444-4444-4444-8444-444444444444")
	cleanup := func() {
		db.Exec(`DELETE FROM galeri WHERE id = ?`, galleryID)
		db.Exec(`DELETE FROM events WHERE id IN (?, ?)`, eventID, draftEventID)
		db.Exec(`DELETE FROM blogs WHERE id IN (?, ?)`, blogID, draftBlogID)
	}
	cleanup()
	t.Cleanup(cleanup)

	if err := db.Exec(`INSERT INTO auth.users (id,email) VALUES (?,?)
		ON CONFLICT (id) DO NOTHING`, creator, "integration@example.com").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO profiles (id,email,username,role) VALUES (?,?,?,'admin')
		ON CONFLICT (id) DO UPDATE SET username=EXCLUDED.username`,
		creator, "integration@example.com", "Integration Admin").Error; err != nil {
		t.Fatal(err)
	}

	blogRepo := blog.NewRepository(db)
	blogItem := &models.Blog{
		ID: blogID, Title: "Published blog",
		Author: "Admin", Category: "integration_test", Content: "Content", Status: "PUBLISHED",
		CreatedBy: &creator,
	}
	if err := blogRepo.Create(ctx, blogItem); err != nil {
		t.Fatal(err)
	}
	draftBlog := &models.Blog{
		ID: draftBlogID, Title: "Draft blog", Author: "Admin", Category: "integration_test",
		Content: "Draft", Status: "DRAFT", CreatedBy: &creator,
	}
	if err := blogRepo.Create(ctx, draftBlog); err != nil {
		t.Fatal(err)
	}
	blogs, total, err := blogRepo.List(ctx, shared.ListOptions{
		Limit: 10, Published: true, Category: "integration_test", Sort: "title_asc",
	})
	if err != nil || total != 1 || len(blogs) != 1 || blogs[0].AuthorProfile == nil {
		t.Fatalf("blog repository contract failed: total=%d items=%d err=%v", total, len(blogs), err)
	}
	foundBlog, err := blogRepo.Find(ctx, string(blogID), true)
	if err != nil || foundBlog.AuthorProfile == nil {
		t.Fatalf("blog find/preload failed: %v", err)
	}
	if _, err := blogRepo.Find(ctx, string(draftBlogID), true); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("public find must hide draft blog, got %v", err)
	}
	foundBlog.Title = "Updated blog"
	if err := blogRepo.Update(ctx, foundBlog); err != nil {
		t.Fatal(err)
	}
	foundBlog, err = blogRepo.Find(ctx, string(blogID), false)
	if err != nil || foundBlog.Title != "Updated blog" {
		t.Fatalf("blog update contract failed: title=%q err=%v", foundBlog.Title, err)
	}
	if err := blogRepo.Delete(ctx, string(draftBlogID)); err != nil {
		t.Fatal(err)
	}
	if _, err := blogRepo.Find(ctx, string(draftBlogID), false); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("deleted blog must be absent, got %v", err)
	}

	eventRepo := event.NewRepository(db)
	eventItem := &models.Event{
		ID: eventID, Title: "Published event",
		Description: "Description", Author: "Admin", Category: "integration_test",
		EventDate: "2026-08-01", Status: "UPCOMING", PublicationStatus: "PUBLISHED",
		CreatedBy: &creator,
	}
	if err := eventRepo.Create(ctx, eventItem); err != nil {
		t.Fatal(err)
	}
	draftEvent := &models.Event{
		ID: draftEventID, Title: "Draft event", Description: "Draft", Author: "Admin",
		Category: "integration_test", EventDate: "2026-08-10", Status: "UPCOMING",
		PublicationStatus: "DRAFT", CreatedBy: &creator,
	}
	if err := eventRepo.Create(ctx, draftEvent); err != nil {
		t.Fatal(err)
	}
	events, total, err := eventRepo.List(ctx, shared.ListOptions{
		Limit: 10, Published: true, Category: "integration_test",
		StartDate: "2026-08-01", EndDate: "2026-08-02", Sort: "oldest",
	})
	if err != nil || total != 1 || len(events) != 1 || events[0].AuthorProfile == nil {
		t.Fatalf("event repository contract failed: total=%d items=%d err=%v", total, len(events), err)
	}
	foundEvent, err := eventRepo.Find(ctx, string(eventID), true)
	if err != nil || foundEvent.AuthorProfile == nil {
		t.Fatalf("event find/preload failed: %v", err)
	}
	if _, err := eventRepo.Find(ctx, string(draftEventID), true); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("public find must hide draft event, got %v", err)
	}
	foundEvent.Status = "ONGOING"
	if err := eventRepo.Update(ctx, foundEvent); err != nil {
		t.Fatal(err)
	}
	foundEvent, err = eventRepo.Find(ctx, string(eventID), false)
	if err != nil || foundEvent.Status != "ONGOING" {
		t.Fatalf("event update contract failed: status=%q err=%v", foundEvent.Status, err)
	}
	if err := eventRepo.Delete(ctx, string(draftEventID)); err != nil {
		t.Fatal(err)
	}
	if _, err := eventRepo.Find(ctx, string(draftEventID), false); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("deleted event must be absent, got %v", err)
	}

	galleryRepo := gallery.NewRepository(db)
	galleryItem := &models.Gallery{
		ID: galleryID, Title: "Gallery",
		Link: "https://example.com/gallery", Category: "teknik_informatika", TakenAt: "2026-08-01",
		CreatedBy: &creator,
	}
	if err := galleryRepo.Create(ctx, galleryItem); err != nil {
		t.Fatal(err)
	}
	galleries, total, err := galleryRepo.List(ctx, shared.ListOptions{
		Limit: 10, Category: "teknik_informatika", Sort: "latest",
	})
	if err != nil || total != 1 || len(galleries) != 1 || galleries[0].AuthorProfile == nil {
		t.Fatalf("gallery repository contract failed: total=%d items=%d err=%v", total, len(galleries), err)
	}
	foundGallery, err := galleryRepo.Find(ctx, string(galleryID))
	if err != nil || foundGallery.AuthorProfile == nil {
		t.Fatalf("gallery find/preload failed: %v", err)
	}
	foundGallery.Title = "Updated gallery"
	if err := galleryRepo.Update(ctx, foundGallery); err != nil {
		t.Fatal(err)
	}
	foundGallery, err = galleryRepo.Find(ctx, string(galleryID))
	if err != nil || foundGallery.Title != "Updated gallery" {
		t.Fatalf("gallery update contract failed: title=%q err=%v", foundGallery.Title, err)
	}
	if err := galleryRepo.Delete(ctx, string(galleryID)); err != nil {
		t.Fatal(err)
	}
	if _, err := galleryRepo.Find(ctx, string(galleryID)); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("deleted gallery must be absent, got %v", err)
	}
}
