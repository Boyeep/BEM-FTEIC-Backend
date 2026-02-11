package models

import (
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	ImagePath   string `json:"image_path"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
	AuthorPhoto string `json:"author_photo"`
}
