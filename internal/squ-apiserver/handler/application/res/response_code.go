package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrApplicationNotFound = 70001
	ErrDeployFailed         = 70002
	ErrAlreadyDeployed      = 70003
)

func RegisterCode() {
	response.Register(ErrApplicationNotFound, "application not found")
	response.Register(ErrDeployFailed, "deploy failed")
	response.Register(ErrAlreadyDeployed, "application already deployed to this server")
}
