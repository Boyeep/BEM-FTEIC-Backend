package dto

type CreateBlog struct {
	Title       string `json:"title" binding:"required,max=200"`
	Excerpt     string `json:"excerpt" binding:"max=500"`
	Author      string `json:"author" binding:"required,max=120"`
	Category    string `json:"category" binding:"required,max=80"`
	CoverImage  string `json:"cover_image" binding:"omitempty,url"`
	Content     string `json:"content" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
	PublishedAt string `json:"published_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateBlog = CreateBlog

type CreateEvent struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"required"`
	Author      string `json:"author" binding:"required,max=120"`
	Category    string `json:"category" binding:"required,max=80"`
	CoverImage  string `json:"cover_image" binding:"omitempty,url"`
	EventDate   string `json:"event_date" binding:"required,datetime=2006-01-02"`
	Status      string `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
}

type UpdateEvent = CreateEvent

type CreateGallery struct {
	Title    string `json:"title" binding:"required,max=200"`
	Link     string `json:"link" binding:"required,url"`
	ImageURL string `json:"image_url" binding:"omitempty,url"`
	TakenAt  string `json:"taken_at" binding:"required,datetime=2006-01-02"`
}

type UpdateGallery = CreateGallery
