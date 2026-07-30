package content

import (
	"errors"

	"repo-backend/internal/dto"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"
	"repo-backend/pkg/apperr"
)

type ListOptions struct {
	Offset    int
	Limit     int
	Category  string
	StartDate string
	EndDate   string
	Sort      string
	Published bool
}

func Options(query dto.ListQuery, published bool) ListOptions {
	return ListOptions{
		Offset:    (query.Page - 1) * query.PageSize,
		Limit:     query.PageSize,
		Category:  query.Category,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
		Sort:      query.Sort,
		Published: published,
	}
}

func Page[T any](items []T, page, pageSize int, total int64) *models.Page[T] {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	return &models.Page[T]{
		Items: items,
		Pagination: models.Pagination{
			Page:            page,
			PageSize:        pageSize,
			TotalItems:      total,
			TotalPages:      totalPages,
			HasNextPage:     page < totalPages,
			HasPreviousPage: page > 1,
		},
	}
}

func ServiceError(err error, resource string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperr.NotFound(resource)
	}
	return apperr.Internal(err)
}

func Order(sort, dateColumn string) string {
	switch sort {
	case "oldest":
		return dateColumn + " ASC"
	case "title_asc":
		return "title ASC"
	case "title_desc":
		return "title DESC"
	default:
		return dateColumn + " DESC"
	}
}
