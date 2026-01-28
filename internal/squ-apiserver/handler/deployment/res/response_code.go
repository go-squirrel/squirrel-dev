package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrApplicationNotFound    = 70001
	ErrDeploymentNotFound     = 70002
	ErrDeployFailed           = 70003
	ErrAlreadyDeployed        = 70004
	ErrApplicationNotDeployed = 70005
)

func RegisterCode() {
	response.Register(ErrApplicationNotFound, "application not found")
	response.Register(ErrDeploymentNotFound, "deployment not found")
	response.Register(ErrDeployFailed, "deploy failed")
	response.Register(ErrAlreadyDeployed, "application already deployed to this server")
	response.Register(ErrApplicationNotDeployed, "application not deployed to this server")
}
