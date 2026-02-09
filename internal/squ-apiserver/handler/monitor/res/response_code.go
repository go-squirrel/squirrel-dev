package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrMonitorFailed        = 81001
	ErrInvalidMonitorConfig = 81002
	ErrMonitorDataNotFound  = 81003
	ErrServerNotFound       = 81004
)

func RegisterCode() {
	response.Register(ErrMonitorFailed, "monitor request failed")
	response.Register(ErrInvalidMonitorConfig, "invalid monitor configuration")
	response.Register(ErrMonitorDataNotFound, "monitor data not found")
	response.Register(ErrServerNotFound, "server not found")
}
