package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrServerNotFound = 80001
	ErrMonitorFailed  = 80002
)

func RegisterCode() {
	response.Register(ErrServerNotFound, "server not found")
	response.Register(ErrMonitorFailed, "monitor request failed")
}
