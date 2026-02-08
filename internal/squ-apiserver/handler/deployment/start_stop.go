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
	// 1. 根据deployment.ID查询部署记录
	deployment, err := a.Repository.Get(deploymentID)
	if err != nil {
		return response.Error(returnDeploymentErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(deployment.ServerID)
	if err != nil {
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 3. 调用 Agent 停止应用，使用deployID
	stopUrl := fmt.Sprintf("application/stop/%d", deployment.DeployID)

	agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		stopUrl)
	respBody, err := a.HTTPClient.Post(agentURL, nil, nil)
	if err != nil {
		zap.L().Error("调用 Agent 停止应用失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentStopFailed)
	}

	// 解析响应，检查是否停止成功
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentResponseParseFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 停止失败",
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
	// 1. 根据deployment.ID查询部署记录
	deployment, err := a.Repository.Get(deploymentID)
	if err != nil {
		return response.Error(returnDeploymentErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(deployment.ServerID)
	if err != nil {
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 3. 调用 Agent 启动应用，使用deployID
	startUrl := fmt.Sprintf("application/start/%d", deployment.DeployID)

	agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		startUrl)
	respBody, err := a.HTTPClient.Post(agentURL, nil, nil)
	if err != nil {
		zap.L().Error("调用 Agent 启动应用失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentStartFailed)
	}

	// 解析响应，检查是否启动成功
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentResponseParseFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 启动失败",
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrAgentOperationFailed)
	}

	return response.Success("start success")
}
