package application

import "errors"

var (
	ErrInvalidCompose  = errors.New("invalid compose content")
	ErrUnsupportedType = errors.New("unsupported application type")
)
