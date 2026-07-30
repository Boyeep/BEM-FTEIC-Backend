package models

import "time"

type Event struct {
	ID                string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Author            string    `json:"author"`
	Category          string    `json:"category"`
	CoverImage        string    `json:"cover_image"`
	EventDate         string    `json:"event_date" gorm:"type:date"`
	Status            string    `json:"status"`
	PublicationStatus string    `json:"publication_status"`
	CreatedBy         *string   `json:"created_by" gorm:"type:uuid"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
