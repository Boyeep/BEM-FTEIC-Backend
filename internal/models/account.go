package models

import "time"

type Profile struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Role      string  `json:"role"`
}

type PublicProfile struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
}

type WhitelistEntry struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedBy *string   `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
