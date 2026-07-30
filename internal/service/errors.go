package service

import (
	"errors"

	"repo-backend/internal/repository"
	"repo-backend/pkg/apperr"
)

func serviceError(err error, resource string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperr.NotFound(resource)
	}
	return apperr.Internal(err)
}
