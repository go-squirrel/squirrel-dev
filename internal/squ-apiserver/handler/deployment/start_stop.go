package deployment

import (
	"encoding/json"
	"fmt"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"
	"squirrel-dev/pkg/utils"

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

	agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		stopUrl)
	respBody, err := a.HTTPClient.Post(agentURL, nil, nil)
	if err != nil {
		zap.L().Error("failed to call agent to stop application",
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("deployment_id", deploymentID),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentStopFailed)
	}

	// Parse response, check if stop succeeded
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("failed to parse agent response",
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("deployment_id", deploymentID),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentResponseParseFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("agent stop application failed",
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("deployment_id", deploymentID),
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrAgentOperationFailed)
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

	agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		startUrl)
	respBody, err := a.HTTPClient.Post(agentURL, nil, nil)
	if err != nil {
		zap.L().Error("failed to call agent to start application",
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("deployment_id", deploymentID),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentStartFailed)
	}

	// Parse response, check if start succeeded
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("failed to parse agent response",
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("deployment_id", deploymentID),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentResponseParseFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("agent start application failed",
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("deployment_id", deploymentID),
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrAgentOperationFailed)
	}

	return response.Success("start success")
}
