package deployment

import (
	"encoding/json"
	"fmt"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/req"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"
	"squirrel-dev/internal/squ-apiserver/model"
	appRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	deploymentRepository "squirrel-dev/internal/squ-apiserver/repository/deployment"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"
	"time"

	"go.uber.org/zap"
)

type Deployment struct {
	Config     *config.Config
	Repository deploymentRepository.Repository
	AppRepo    appRepository.Repository
	ServerRepo serverRepository.Repository
	HTTPClient *httpclient.Client
}

func New(config *config.Config, repo deploymentRepository.Repository, appRepo appRepository.Repository, serverRepo serverRepository.Repository) *Deployment {
	hc := httpclient.NewClient(30 * time.Second)
	return &Deployment{
		Config:     config,
		Repository: repo,
		AppRepo:    appRepo,
		ServerRepo: serverRepo,
		HTTPClient: hc,
	}
}

func (a *Deployment) Deploy(request req.DeployApplication) response.Response {
	// 1. 检查应用是否存在
	app, err := a.AppRepo.Get(request.ApplicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(request.ServerID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	deployID, err := utils.IDGenerate()
	if err != nil {
		zap.L().Error("生成部署ID失败",
			zap.Error(err),
		)
		return response.Error(res.ErrDeployIDGenerateFailed)
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
		DeployID:    deployID,
	}
	respBody, err := a.HTTPClient.Post(agentURL, agentReq, nil)
	if err != nil {
		zap.L().Error("部署请求发送失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentRequestFailed)
	}

	// 解析响应，检查是否部署成功
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrAgentResponseParseFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 部署失败",
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		return response.Error(res.ErrAgentDeployFailed)
	}

	// 5. 创建应用服务器关联记录
	appServer := model.Deployment{
		ServerID:      request.ServerID,
		ApplicationID: request.ApplicationID,
		DeployID:      deployID,
	}

	err = a.Repository.Add(&appServer)
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
func (a *Deployment) ListServers(applicationID uint) response.Response {
	// 检查应用是否存在
	_, err := a.AppRepo.Get(applicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 查询应用服务器关联记录
	appServers, err := a.Repository.List(applicationID)
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
			ID:        server.ID,
			IpAddress: server.IpAddress,
			AgentPort: server.AgentPort,
		})
	}

	return response.Success(servers)
}

// Undeploy 取消部署应用
func (a *Deployment) Undeploy(deploymentID uint) response.Response {
	// 1. 根据deployment.ID查询部署记录
	deployment, err := a.Repository.Get(deploymentID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrDeploymentNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(deployment.ServerID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 调用 Agent 删除应用，使用deployID
	deleteUrl := fmt.Sprintf("application/delete/%d", deployment.DeployID)

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
		return response.Error(res.ErrAgentDeleteFailed)
	}

	// 4. 删除应用服务器关联记录
	err = a.Repository.Delete(deployment.ID)
	if err != nil {
		zap.L().Error("删除应用服务器关联记录失败",
			zap.Uint("id", deployment.ID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("undeploy success")
}

func (a *Deployment) ReportStatus(request req.ReportApplicationStatus) response.Response {
	// 验证应用服务器关联记录是否存在，使用deployID
	_, err := a.Repository.GetByDeployID(request.DeployID)
	if err != nil {
		zap.L().Error("应用服务器关联记录不存在",
			zap.Uint64("deploy_id", request.DeployID),
			zap.Error(err),
		)
		return response.Error(response.ErrCodeParameter)
	}

	// 更新状态，使用deployID
	err = a.Repository.UpdateStatus(request.DeployID, request.Status)
	if err != nil {
		zap.L().Error("更新应用状态失败",
			zap.Uint64("deploy_id", request.DeployID),
			zap.String("status", request.Status),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	zap.L().Info("应用状态已更新",
		zap.Uint64("deploy_id", request.DeployID),
		zap.String("status", request.Status),
	)

	return response.Success("success")
}

// List 列出部署应用
func (a *Deployment) List(serverID uint) response.Response {
	// 查询部署列表
	deployments, err := a.Repository.List(serverID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 构建响应数据
	var result []res.Deployment
	for _, deployment := range deployments {
		// 获取应用信息
		app, err := a.AppRepo.Get(deployment.ApplicationID)
		if err != nil {
			zap.L().Warn("获取应用信息失败",
				zap.Uint("application_id", deployment.ApplicationID),
				zap.Error(err),
			)
			continue
		}

		// 获取服务器信息
		server, err := a.ServerRepo.Get(deployment.ServerID)
		if err != nil {
			zap.L().Warn("获取服务器信息失败",
				zap.Uint("server_id", deployment.ServerID),
				zap.Error(err),
			)
			continue
		}

		result = append(result, res.Deployment{
			ID:       deployment.ID,
			DeployID: deployment.DeployID,
			Application: res.ApplicationInfo{
				ID:          app.ID,
				Name:        app.Name,
				Description: app.Description,
				Type:        app.Type,
				Version:     app.Version,
			},
			Server: res.ServerInfo{
				ID:        server.ID,
				IpAddress: server.IpAddress,
				AgentPort: server.AgentPort,
			},
			Status:     deployment.Status,
			DeployedAt: deployment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.Success(result)
}
