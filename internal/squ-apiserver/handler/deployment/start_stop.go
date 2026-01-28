package deployment

import (
	"encoding/json"
	"fmt"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

// Stop 停止应用
func (a *Deployment) Stop(applicationID, serverID uint) response.Response {
	// 1. 检查应用是否存在
	app, err := a.AppRepo.Get(applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(serverID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 检查是否已部署
	_, err = a.Repository.GetByServerAndApp(serverID, applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotDeployed)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 4. 调用 Agent 停止应用
	stopUrl := fmt.Sprintf("application/stop/%s", app.Name)

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
		return response.Error(res.ErrDeployFailed)
	}

	// 解析响应，检查是否停止成功
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 停止失败",
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrDeployFailed)
	}

	return response.Success("stop success")
}

// Start 启动应用
func (a *Deployment) Start(applicationID, serverID uint) response.Response {
	// 1. 检查应用是否存在
	app, err := a.AppRepo.Get(applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(serverID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 检查是否已部署
	_, err = a.Repository.GetByServerAndApp(serverID, applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotDeployed)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 4. 调用 Agent 启动应用
	stopUrl := fmt.Sprintf("application/start/%s", app.Name)

	agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		stopUrl)
	respBody, err := a.HTTPClient.Post(agentURL, nil, nil)
	if err != nil {
		zap.L().Error("调用 Agent 启动应用失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployFailed)
	}

	// 解析响应，检查是否启动成功
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 启动失败",
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrDeployFailed)
	}

	return response.Success("start success")
}
