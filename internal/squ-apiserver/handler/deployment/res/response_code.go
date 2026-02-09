package res

import "squirrel-dev/internal/pkg/response"

const (
	// 72000-72019: 基础操作
	ErrDeploymentNotFound           = 72001
	ErrAlreadyDeployed              = 72002
	ErrApplicationNotDeployed       = 72003
	ErrDeployIDGenerateFailed       = 72004
	ErrCreateDeploymentRecordFailed = 72005
	ErrInvalidDeploymentConfig      = 72006

	// 72020-72039: Agent 相关
	ErrAgentRequestFailed       = 72021
	ErrAgentResponseParseFailed = 72022
	ErrAgentDeployFailed        = 72023
	ErrAgentDeleteFailed        = 72024
	ErrAgentStopFailed          = 72025
	ErrAgentStartFailed         = 72026
	ErrAgentOperationFailed     = 72027
)

func RegisterCode() {
	// 72000-72019: 基础操作
	response.Register(ErrDeploymentNotFound, "deployment not found")
	response.Register(ErrAlreadyDeployed, "application already deployed to this server")
	response.Register(ErrApplicationNotDeployed, "application not deployed to this server")
	response.Register(ErrDeployIDGenerateFailed, "failed to generate deploy ID")
	response.Register(ErrCreateDeploymentRecordFailed, "failed to create deployment record")
	response.Register(ErrInvalidDeploymentConfig, "invalid deployment configuration")

	// 72020-72039: Agent 相关
	response.Register(ErrAgentRequestFailed, "failed to send request to agent")
	response.Register(ErrAgentResponseParseFailed, "failed to parse agent response")
	response.Register(ErrAgentDeployFailed, "agent deployment failed")
	response.Register(ErrAgentDeleteFailed, "agent delete application failed")
	response.Register(ErrAgentStopFailed, "agent stop application failed")
	response.Register(ErrAgentStartFailed, "agent start application failed")
	response.Register(ErrAgentOperationFailed, "agent operation failed")
}
