package dto

type CreateBlog struct {
	Title       string `json:"title" binding:"required,max=200"`
	Excerpt     string `json:"excerpt" binding:"max=500"`
	Author      string `json:"author" binding:"required,max=120"`
	Category    string `json:"category" binding:"required,max=80"`
	CoverImage  string `json:"cover_image" binding:"omitempty,url"`
	Content     string `json:"content" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	PublishedAt string `json:"published_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateBlog = CreateBlog

type CreateEvent struct {
	Title             string `json:"title" binding:"required,max=200"`
	Description       string `json:"description" binding:"required"`
	Author            string `json:"author" binding:"required,max=120"`
	Category          string `json:"category" binding:"required,max=80"`
	CoverImage        string `json:"cover_image" binding:"omitempty,url"`
	EventDate         string `json:"event_date" binding:"required,datetime=2006-01-02"`
	Status            string `json:"status" binding:"required,oneof=UPCOMING ONGOING ENDED"`
	PublicationStatus string `json:"publication_status" binding:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
}

type UpdateEvent = CreateEvent

type CreateGallery struct {
	Title    string `json:"title" binding:"required,max=200"`
	Link     string `json:"link" binding:"required,url"`
	ImageURL string `json:"image_url" binding:"omitempty,url"`
	Category string `json:"category" binding:"omitempty,oneof=all teknik_elektro teknik_informatika sistem_informasi teknik_komputer teknik_biomedik teknologi_informasi"`
	TakenAt  string `json:"taken_at" binding:"required,datetime=2006-01-02"`
}

type UpdateGallery = CreateGallery

type ListQuery struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Category  string `form:"category" binding:"omitempty,max=80"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Sort      string `form:"sort" binding:"omitempty,oneof=latest oldest title_asc title_desc"`
}

func (q *ListQuery) Defaults(defaultPageSize int) {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = defaultPageSize
	}
}
