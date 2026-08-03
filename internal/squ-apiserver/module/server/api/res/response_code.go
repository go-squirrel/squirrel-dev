package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrServerNotFound      = 60001
	ErrServerAliasExists   = 60002
	ErrServerUUIDNotFound  = 60003
	ErrServerAlreadyExists = 60004
	ErrServerUpdateFailed  = 60005
	ErrServerDeleteFailed  = 60006

	ErrInvalidParameter = 60021
	ErrInvalidAuthType  = 60022
	ErrInvalidSSHConfig = 60023
	ErrSSHTestFailed    = 60024

	ErrConnectFailed      = 60041
	ErrAgentOffline       = 60042
	ErrAgentRequestFailed = 60043
)

func RegisterCode() {
	response.Register(ErrServerNotFound, "server not found")
	response.Register(ErrServerAliasExists, "server alias already exists")
	response.Register(ErrServerUUIDNotFound, "server not found by UUID")
	response.Register(ErrServerAlreadyExists, "server already exists")
	response.Register(ErrServerUpdateFailed, "server update failed")
	response.Register(ErrServerDeleteFailed, "server delete failed")

	response.Register(ErrInvalidParameter, "invalid parameter")
	response.Register(ErrInvalidAuthType, "invalid auth type")
	response.Register(ErrInvalidSSHConfig, "invalid SSH configuration")
	response.Register(ErrSSHTestFailed, "SSH connection test failed")

	response.Register(ErrConnectFailed, "connect failed")
	response.Register(ErrAgentOffline, "agent is offline")
	response.Register(ErrAgentRequestFailed, "agent request failed")
}
