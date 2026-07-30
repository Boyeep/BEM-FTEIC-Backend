package database

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		version := entry.Name()
		var count int64
		if err := db.Raw(
			"SELECT COUNT(*) FROM schema_migrations WHERE version = ?",
			version,
		).Scan(&count).Error; err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if count > 0 {
			continue
		}

		sql, err := migrationFiles.ReadFile("migrations/" + version)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", version, err)
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(sql)).Error; err != nil {
				return err
			}
			return tx.Exec(
				"INSERT INTO schema_migrations (version) VALUES (?)",
				version,
			).Error
		}); err != nil {
			return fmt.Errorf("apply migration %s: %w", version, err)
		}
	}

	return nil
}
