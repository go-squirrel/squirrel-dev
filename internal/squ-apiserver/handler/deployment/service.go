package deployment

import (
	"context"
	"fmt"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/agent"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/req"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"
	"squirrel-dev/internal/squ-apiserver/model"
	appRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	deploymentRepository "squirrel-dev/internal/squ-apiserver/repository/deployment"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

type Deployment struct {
	Config      *config.Config
	Repository  deploymentRepository.Repository
	AppRepo     appRepository.Repository
	ServerRepo  serverRepository.Repository
	AgentClient *agent.Client
}

func New(config *config.Config, repo deploymentRepository.Repository, appRepo appRepository.Repository, serverRepo serverRepository.Repository) *Deployment {
	return &Deployment{
		Config:      config,
		Repository:  repo,
		AppRepo:     appRepo,
		ServerRepo:  serverRepo,
		AgentClient: agent.NewClient(config),
	}
}

func (a *Deployment) Deploy(request req.DeployApplication) response.Response {
	// 1. Check if application exists
	app, err := a.AppRepo.Get(request.ApplicationID)
	if err != nil {
		zap.L().Error("failed to get application for deployment", zap.Uint("application_id", request.ApplicationID), zap.Uint("server_id", request.ServerID), zap.Error(err))
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 2. Check if server exists
	server, err := a.ServerRepo.Get(request.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for deployment",
			zap.Uint("application_id", request.ApplicationID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 3. Check compose content conflicts with existing deployments on this server
	existingDeployments, err := a.Repository.List(request.ServerID)
	if err != nil {
		zap.L().Error("failed to list existing deployments for conflict check",
			zap.Uint("application_id", request.ApplicationID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	for _, deployment := range existingDeployments {

		// 获取已部署应用的信息
		existingApp, err := a.AppRepo.Get(deployment.ApplicationID)
		if err != nil {
			zap.L().Warn("failed to get existing application for conflict check",
				zap.Uint("application_id", deployment.ApplicationID),
				zap.Uint("server_id", request.ServerID),
				zap.Error(err),
			)
			continue
		}

		// 检查 compose 内容冲突
		if existingApp.Type == "compose" && app.Type == "compose" {
			conflictCode := checkComposeContent(app.Content, deployment.Content)
			if conflictCode != 0 {
				zap.L().Warn("compose conflict detected",
					zap.Uint("application_id", request.ApplicationID),
					zap.String("application_name", app.Name),
					zap.Uint("existing_application_id", deployment.ID),
					zap.String("existing_application_name", existingApp.Name),
					zap.Uint("server_id", request.ServerID),
					zap.Int("conflict_code", conflictCode),
				)
				return response.Error(conflictCode)
			}
		}
	}

	deployID, err := utils.IDGenerate()
	if err != nil {
		zap.L().Error("failed to generate deployment ID",
			zap.Uint("application_id", request.ApplicationID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrDeployIDGenerateFailed)
	}

	// 5. Send deployment request to agent
	agentReq := req.ApplicationAgent{
		Name:        app.Name,
		Description: app.Description,
		Type:        app.Type,
		Content:     app.Content,
		Version:     app.Version,
		ServerID:    request.ServerID,
		DeployID:    deployID,
	}
	result := a.AgentClient.Post(context.Background(), server, "application", agentReq,
		zap.Uint64("deploy_id", deployID),
		zap.Uint("application_id", request.ApplicationID),
		zap.Uint("server_id", request.ServerID),
	)
	if result.Err != nil {
		return response.Error(res.ErrAgentDeployFailed)
	}

	// 6. Create application server association record
	appServer := model.Deployment{
		ServerID:      request.ServerID,
		ApplicationID: request.ApplicationID,
		Content:       app.Content,
		DeployID:      deployID,
	}

	err = a.Repository.Add(&appServer)
	if err != nil {
		zap.L().Error("failed to create application server association record",
			zap.Uint64("deploy_id", deployID),
			zap.Uint("application_id", request.ApplicationID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)

		// Rollback: attempt to delete deployed application on agent
		deleteUrl := fmt.Sprintf("application/delete/%s", app.Name)
		zap.L().Info("rollback: attempting to delete deployed application on agent",
			zap.Uint64("deploy_id", deployID),
			zap.Uint("application_id", request.ApplicationID),
			zap.Uint("server_id", request.ServerID),
		)
		_ = a.AgentClient.Post(context.Background(), server, deleteUrl, nil,
			zap.Uint64("deploy_id", deployID),
			zap.Uint("application_id", request.ApplicationID),
			zap.Uint("server_id", request.ServerID),
			zap.String("operation", "rollback"),
		)

		return response.Error(returnDeploymentErrCode(err))
	}

	return response.Success("deploy success")
}

// ListServers 查询应用部署的服务器列表
func (a *Deployment) ListServers(applicationID uint) response.Response {
	// Check if application exists
	_, err := a.AppRepo.Get(applicationID)
	if err != nil {
		zap.L().Error("failed to get application for listing servers",
			zap.Uint("application_id", applicationID),
			zap.Error(err),
		)
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// Query application server association records
	appServers, err := a.Repository.List(applicationID)
	if err != nil {
		zap.L().Error("failed to list application servers",
			zap.Uint("application_id", applicationID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	var servers []res.ServerInfo
	for _, appServer := range appServers {
		server, err := a.ServerRepo.Get(appServer.ServerID)
		if err != nil {
			zap.L().Warn("failed to get server information",
				zap.Uint("application_id", applicationID),
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
	// 1. Query deployment record by deployment ID
	deployment, err := a.Repository.Get(deploymentID)
	if err != nil {
		zap.L().Error("failed to get deployment for undeploy",
			zap.Uint("deployment_id", deploymentID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	// 2. Check if server exists
	server, err := a.ServerRepo.Get(deployment.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for undeploy",
			zap.Uint("deployment_id", deploymentID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrApplicationNotDeployed)
	}

	// 3. Call agent to delete application, use deployID
	deleteUrl := fmt.Sprintf("application/delete/%d", deployment.DeployID)
	result := a.AgentClient.Post(context.Background(), server, deleteUrl, nil,
		zap.Uint64("deploy_id", deployment.DeployID),
		zap.Uint("deployment_id", deploymentID),
	)
	if result.Err != nil {
		return response.Error(res.ErrAgentDeleteFailed)
	}

	// 4. Delete application server association record
	err = a.Repository.Delete(deployment.ID)
	if err != nil {
		zap.L().Error("failed to delete application server association record",
			zap.Uint("deployment_id", deploymentID),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	return response.Success("undeploy success")
}

func (a *Deployment) ReportStatus(request req.ReportApplicationStatus) response.Response {
	// Validate application server association record exists, use deployID
	_, err := a.Repository.GetByDeployID(request.DeployID)
	if err != nil {
		zap.L().Error("application server association record not found",
			zap.Uint64("deploy_id", request.DeployID),
			zap.String("status", request.Status),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	// Update status, use deployID
	err = a.Repository.UpdateStatus(request.DeployID, request.Status)
	if err != nil {
		zap.L().Error("failed to update application status",
			zap.Uint64("deploy_id", request.DeployID),
			zap.String("status", request.Status),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	zap.L().Info("application status updated",
		zap.Uint64("deploy_id", request.DeployID),
		zap.String("status", request.Status),
	)

	return response.Success("success")
}

// List 列出部署应用
func (a *Deployment) List(serverID uint) response.Response {
	// Query deployment list
	deployments, err := a.Repository.List(serverID)
	if err != nil {
		zap.L().Error("failed to list deployments",
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		return response.Error(returnDeploymentErrCode(err))
	}

	// Build response data
	var result []res.Deployment
	for _, deployment := range deployments {
		// Get application information
		app, err := a.AppRepo.Get(deployment.ApplicationID)
		if err != nil {
			zap.L().Warn("failed to get application information",
				zap.Uint("server_id", serverID),
				zap.Uint("application_id", deployment.ApplicationID),
				zap.Uint64("deploy_id", deployment.DeployID),
				zap.Error(err),
			)
			continue
		}

		// Get server information
		server, err := a.ServerRepo.Get(deployment.ServerID)
		if err != nil {
			zap.L().Warn("failed to get server information",
				zap.Uint("server_id", serverID),
				zap.Uint("server_id", deployment.ServerID),
				zap.Uint64("deploy_id", deployment.DeployID),
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
			Content:    deployment.Content,
			DeployedAt: deployment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.Success(result)
}
