package models

import "time"

type Admin struct {
	ID        uint      `json:"id_user" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
