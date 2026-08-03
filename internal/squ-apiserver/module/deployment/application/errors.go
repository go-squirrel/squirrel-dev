package application

import "errors"

var (
	ErrNotFound           = errors.New("deployment not found")
	ErrAlreadyDeployed    = errors.New("application already deployed")
	ErrApplicationMissing = errors.New("application or server not found")
	ErrIDGeneration       = errors.New("failed to generate deploy ID")
	ErrCreateRecord       = errors.New("failed to create deployment record")
	ErrInvalidConfig      = errors.New("invalid deployment configuration")
	ErrContainerConflict  = errors.New("compose container name conflict")
	ErrPortConflict       = errors.New("compose port conflict")
	ErrVolumeConflict     = errors.New("compose volume conflict")
	ErrNetworkConflict    = errors.New("compose network conflict")
	ErrAgentDeploy        = errors.New("agent deployment failed")
	ErrAgentDelete        = errors.New("agent delete failed")
	ErrAgentStop          = errors.New("agent stop failed")
	ErrAgentStart         = errors.New("agent start failed")
)
