package application

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound       = errors.New("script not found")
	ErrDuplicate      = errors.New("script already exists")
	ErrInvalidContent = errors.New("invalid script content")
	ErrExecuteFailed  = errors.New("script execution failed")
	ErrServerNotFound = errors.New("server not found")
)

func repositoryError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDuplicate
	default:
		return ErrInvalidContent
	}
}
