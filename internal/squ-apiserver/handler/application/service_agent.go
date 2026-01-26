package application

import (
	"encoding/json"
	"fmt"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

// Deploy 部署应用到指定服务器
func (a *Application) Deploy(request req.DeployApplication) response.Response {
	// 1. 检查应用是否存在
	app, err := a.Repository.Get(request.ApplicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(request.ServerID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 检查是否已经部署到该服务器
	_, err = a.AppServerRepo.GetByServerAndApp(request.ServerID, request.ApplicationID)
	if err == nil {
		return response.Error(res.ErrAlreadyDeployed)
	}

	// 4. 发送部署请求到 agent
	agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		"application")
	agentReq := req.ApplicationAgent{
		Name:        app.Name,
		Description: app.Description,
		Type:        app.Type,
		Content:     app.Content,
		Version:     app.Version,
		ServerID:    request.ServerID,
	}
	respBody, err := a.HTTPClient.Post(agentURL, agentReq, nil)
	if err != nil {
		zap.L().Error("部署请求发送失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployFailed)
	}

	// 解析响应，检查是否部署成功
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 部署失败",
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrDeployFailed)
	}

	// 5. 创建应用服务器关联记录
	appServer := model.ApplicationServer{
		ServerID:      request.ServerID,
		ApplicationID: request.ApplicationID,
	}

	err = a.AppServerRepo.Add(&appServer)
	if err != nil {
		zap.L().Error("创建应用服务器关联记录失败",
			zap.Uint("server_id", request.ServerID),
			zap.Uint("application_id", request.ApplicationID),
			zap.Error(err),
		)

		// 回滚：尝试调用 Agent 删除已部署的应用
		deleteUrl := fmt.Sprintf("application/delete/%s", app.Name)

		agentDeleteURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
			server.IpAddress,
			server.AgentPort,
			a.Config.Agent.Http.BaseUrl,
			deleteUrl)

		zap.L().Info("回滚：尝试删除 Agent 端已部署的应用",
			zap.String("url", agentDeleteURL),
		)
		_, err = a.HTTPClient.Post(agentDeleteURL, nil, nil)
		if err != nil {
			zap.L().Error("回滚失败：删除 Agent 端应用失败",
				zap.String("url", agentDeleteURL),
				zap.Error(err),
			)
		}

		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("deploy success")
}

// ListServers 查询应用部署的服务器列表
func (a *Application) ListServers(applicationID uint) response.Response {
	// 检查应用是否存在
	_, err := a.Repository.Get(applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 查询应用服务器关联记录
	appServers, err := a.AppServerRepo.List(applicationID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	var servers []res.ServerInfo
	for _, appServer := range appServers {
		server, err := a.ServerRepo.Get(appServer.ServerID)
		if err != nil {
			zap.L().Warn("获取服务器信息失败",
				zap.Uint("server_id", appServer.ServerID),
				zap.Error(err),
			)
			continue
		}

		servers = append(servers, res.ServerInfo{
			ServerID:   server.ID,
			IpAddress:  server.IpAddress,
			AgentPort:  server.AgentPort,
			DeployedAt: appServer.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.Success(servers)
}

// Undeploy 取消部署应用
func (a *Application) Undeploy(applicationID, serverID uint) response.Response {
	// 1. 检查应用是否存在
	app, err := a.Repository.Get(applicationID)
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
	appServer, err := a.AppServerRepo.GetByServerAndApp(serverID, applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotDeployed)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 4. 调用 Agent 删除应用

	deleteUrl := fmt.Sprintf("application/delete/%s", app.Name)

	agentDeleteURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		a.Config.Agent.Http.BaseUrl,
		deleteUrl)
	_, err = a.HTTPClient.Post(agentDeleteURL, nil, nil)
	if err != nil {
		zap.L().Error("调用 Agent 删除应用失败",
			zap.String("url", agentDeleteURL),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployFailed)
	}

	// 5. 删除应用服务器关联记录
	err = a.AppServerRepo.Delete(appServer.ID)
	if err != nil {
		zap.L().Error("删除应用服务器关联记录失败",
			zap.Uint("id", appServer.ID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("undeploy success")
}

// Stop 停止应用
func (a *Application) Stop(applicationID, serverID uint) response.Response {
	// 1. 检查应用是否存在
	app, err := a.Repository.Get(applicationID)
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
	_, err = a.AppServerRepo.GetByServerAndApp(serverID, applicationID)
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
func (a *Application) Start(applicationID, serverID uint) response.Response {
	// 1. 检查应用是否存在
	app, err := a.Repository.Get(applicationID)
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
	_, err = a.AppServerRepo.GetByServerAndApp(serverID, applicationID)
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
