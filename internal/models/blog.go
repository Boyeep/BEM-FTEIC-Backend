package models

import (
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description" gorm:"type:text"`
	AdminID     uint   `json:"admin_id"`
	AuthorName  string `json:"author_name"`
	AuthorPhoto string `json:"author_photo"`
}
