package models

import "time"

// Organizer mewakili pilihan penyelenggara event
type Organizer string

const (
	OrganizerFTEIC              Organizer = "FTEIC"
	OrganizerTeknikElektro      Organizer = "Teknik Elektro"
	OrganizerTeknikInformatika  Organizer = "Teknik Informatika"
	OrganizerSistemInformasi    Organizer = "Sistem Informasi"
	OrganizerTeknikKomputer     Organizer = "Teknik Komputer"
	OrganizerTeknikBiomedik     Organizer = "Teknik Biomedik"
	OrganizerTeknologiInformasi Organizer = "Teknologi Informasi"
)

type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Photo       string    `json:"photo"`
	Organizer   Organizer `json:"organizer" gorm:"type:varchar(50);not null"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	IsPublished bool      `json:"is_published" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
