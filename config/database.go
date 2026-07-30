package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDatabase() (*gorm.DB, error) {
	databaseURL, err := Required("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	var db *gorm.DB
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  databaseURL,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("connect to database after 10 attempts: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	return db, nil
}
