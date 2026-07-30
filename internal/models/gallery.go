package models

import "time"

type Gallery struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title     string    `json:"title"`
	Link      string    `json:"link"`
	ImageURL  string    `json:"image_url"`
	Category  string    `json:"category"`
	TakenAt   string    `json:"taken_at" gorm:"type:date"`
	CreatedBy *string   `json:"created_by" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Gallery) TableName() string {
	return "galeri"
}
