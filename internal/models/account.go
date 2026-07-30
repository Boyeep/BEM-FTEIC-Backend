package models

import "time"

type Profile struct {
	ID        UUID    `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Role      string  `json:"role"`
}

type PublicProfile struct {
	ID        UUID    `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
}

func (PublicProfile) TableName() string { return "profiles" }

type WhitelistEntry struct {
	ID        UUID      `json:"id"`
	Email     string    `json:"email"`
	CreatedBy *UUID     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
