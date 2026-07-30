package database

import "gorm.io/gorm"

// DB is kept as a compatibility bridge for the current handlers. New modules
// should receive *gorm.DB through bootstrap instead of reading this global.
var DB *gorm.DB

func Set(db *gorm.DB) {
	DB = db
}
