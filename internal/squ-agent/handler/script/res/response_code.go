package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrScriptExecutionFailed = 90001
	ErrScriptAlreadyRunning  = 90002
)

func RegisterCode() {
	response.Register(ErrScriptExecutionFailed, "script execution failed")
	response.Register(ErrScriptAlreadyRunning, "script is already running")
}
