package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrServerNotFound    = 60001
	ErrServerAliasExists = 60002
	ErrInvalidSSHConfig  = 60003
	ErrConnectFailed     = 60004
)

func RegisterCode() {
	response.Register(ErrServerNotFound, "server not found")
	response.Register(ErrServerAliasExists, "server alias already exists")
	response.Register(ErrInvalidSSHConfig, "invalid SSH configuration")
	response.Register(ErrConnectFailed, "connect failed")
}
