package repository

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("repository: not found")
var ErrConflict = errors.New("repository: conflict")
var ErrLastAdmin = errors.New("repository: cannot remove the last admin")

func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
