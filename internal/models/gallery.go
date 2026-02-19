package models

import "time"

type Gallery struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	DriveLink string    `json:"drive_link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
