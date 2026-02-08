package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrApplicationNotFound       = 70001
	ErrDeploymentNotFound        = 70002
	ErrDeployFailed              = 70003
	ErrAlreadyDeployed           = 70004
	ErrApplicationNotDeployed    = 70005
	ErrDeployIDGenerateFailed    = 70006
	ErrAgentRequestFailed        = 70007
	ErrAgentResponseParseFailed  = 70008
	ErrAgentDeployFailed         = 70009
	ErrAgentDeleteFailed         = 70010
	ErrAgentStopFailed           = 70011
	ErrAgentStartFailed          = 70012
	ErrAgentOperationFailed      = 70013
	ErrServerNotFound            = 70014
	ErrCreateDeploymentRecordFailed = 70015
)

func RegisterCode() {
	response.Register(ErrApplicationNotFound, "application not found")
	response.Register(ErrDeploymentNotFound, "deployment not found")
	response.Register(ErrDeployFailed, "deploy failed")
	response.Register(ErrAlreadyDeployed, "application already deployed to this server")
	response.Register(ErrApplicationNotDeployed, "application not deployed to this server")
	response.Register(ErrDeployIDGenerateFailed, "failed to generate deploy ID")
	response.Register(ErrAgentRequestFailed, "failed to send request to agent")
	response.Register(ErrAgentResponseParseFailed, "failed to parse agent response")
	response.Register(ErrAgentDeployFailed, "agent deployment failed")
	response.Register(ErrAgentDeleteFailed, "agent delete application failed")
	response.Register(ErrAgentStopFailed, "agent stop application failed")
	response.Register(ErrAgentStartFailed, "agent start application failed")
	response.Register(ErrAgentOperationFailed, "agent operation failed")
	response.Register(ErrServerNotFound, "server not found")
	response.Register(ErrCreateDeploymentRecordFailed, "failed to create deployment record")
}
