package models

type Page[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page            int   `json:"page"`
	PageSize        int   `json:"page_size"`
	TotalItems      int64 `json:"total_items"`
	TotalPages      int   `json:"total_pages"`
	HasNextPage     bool  `json:"has_next_page"`
	HasPreviousPage bool  `json:"has_previous_page"`
}
