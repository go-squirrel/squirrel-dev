package application

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrTokenGeneration    = errors.New("failed to generate token")
)
