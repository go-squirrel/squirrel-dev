package application

import "errors"

var (
	ErrMonitorFailed  = errors.New("monitor request failed")
	ErrServerNotFound = errors.New("server not found")
)
