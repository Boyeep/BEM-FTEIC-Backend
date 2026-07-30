package models

import "time"

type Blog struct {
	ID            UUID           `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title         string         `json:"title"`
	Excerpt       string         `json:"excerpt"`
	Author        string         `json:"author"`
	Category      string         `json:"category"`
	CoverImage    string         `json:"cover_image"`
	Content       string         `json:"content"`
	Status        string         `json:"status"`
	PublishedAt   time.Time      `json:"published_at"`
	CreatedBy     *UUID          `json:"created_by" gorm:"type:uuid"`
	AuthorProfile *PublicProfile `json:"author_profile,omitempty" gorm:"foreignKey:CreatedBy;references:ID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
