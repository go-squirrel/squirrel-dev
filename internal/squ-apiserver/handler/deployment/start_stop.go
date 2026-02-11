package deployment

import (
	"context"
	"fmt"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"

	"go.uber.org/zap"
)

// Stop 停止应用
func (a *Deployment) Stop(deploymentID uint) response.Response {
	// 1. Query deployment record by deployment ID
	deployment, err := a.Repository.Get(deploymentID)
	if err != nil {
		zap.L().Error("failed to get deployment for stop",
			zap.Uint("deployment_id", deploymentID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	// 2. Check if server exists
	server, err := a.ServerRepo.Get(deployment.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for stop",
			zap.Uint("deployment_id", deploymentID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 3. Call agent to stop application, use deployID
	stopUrl := fmt.Sprintf("application/stop/%d", deployment.DeployID)

	result := a.AgentClient.Post(context.Background(), server, stopUrl, nil,
		zap.Uint64("deploy_id", deployment.DeployID),
		zap.Uint("deployment_id", deploymentID),
	)
	if result.Err != nil {
		return response.Error(res.ErrAgentStopFailed)
	}

	return response.Success("stop success")
}

// Start 启动应用
func (a *Deployment) Start(deploymentID uint) response.Response {
	// 1. Query deployment record by deployment ID
	deployment, err := a.Repository.Get(deploymentID)
	if err != nil {
		zap.L().Error("failed to get deployment for start",
			zap.Uint("deployment_id", deploymentID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	// 2. Check if server exists
	server, err := a.ServerRepo.Get(deployment.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for start",
			zap.Uint("deployment_id", deploymentID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 3. Call agent to start application, use deployID
	startUrl := fmt.Sprintf("application/start/%d", deployment.DeployID)

	result := a.AgentClient.Post(context.Background(), server, startUrl, nil,
		zap.Uint64("deploy_id", deployment.DeployID),
		zap.Uint("deployment_id", deploymentID),
	)
	if result.Err != nil {
		return response.Error(res.ErrAgentStartFailed)
	}

	return response.Success("start success")
}
