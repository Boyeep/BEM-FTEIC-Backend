package dto

type UpdateProfile struct {
	Username  string `json:"username" binding:"required,min=2,max=80"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

type AddWhitelist struct {
	Email string `json:"email" binding:"required,email,max=254"`
}
