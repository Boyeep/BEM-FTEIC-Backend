package dto

type TrackVisitor struct {
	ID        string `json:"id" binding:"required,max=128"`
	Path      string `json:"path" binding:"required,max=500"`
	UserAgent string `json:"user_agent" binding:"omitempty,max=500"`
}
