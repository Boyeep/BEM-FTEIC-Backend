package models

import "time"

type Event struct {
	ID                UUID           `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Author            string         `json:"author"`
	Category          string         `json:"category"`
	CoverImage        string         `json:"cover_image"`
	EventDate         Date           `json:"event_date" gorm:"type:date"`
	Status            string         `json:"status"`
	PublicationStatus string         `json:"publication_status"`
	CreatedBy         *UUID          `json:"created_by" gorm:"type:uuid"`
	AuthorProfile     *PublicProfile `json:"author_profile,omitempty" gorm:"foreignKey:CreatedBy;references:ID"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
