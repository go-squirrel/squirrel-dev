package application

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/deployment/domain"
)

type AgentApplication struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Version     string `json:"version"`
	ServerID    uint   `json:"server_id"`
	DeployID    uint64 `json:"deploy_id"`
}

type DeploymentView struct {
	Deployment  domain.Deployment
	Application domain.Application
	Server      domain.Server
}

type Service struct {
	repository   domain.Repository
	applications domain.ApplicationReader
	servers      domain.ServerReader
	agent        domain.AgentClient
	ids          domain.IDGenerator
}

func NewService(
	repository domain.Repository,
	applications domain.ApplicationReader,
	servers domain.ServerReader,
	agent domain.AgentClient,
	ids domain.IDGenerator,
) *Service {
	return &Service{repository: repository, applications: applications, servers: servers, agent: agent, ids: ids}
}

func (s *Service) Deploy(ctx context.Context, applicationID, serverID uint) (string, error) {
	app, err := s.applications.Get(ctx, applicationID)
	if err != nil {
		zap.L().Error("failed to get application for deployment",
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		return "", ErrApplicationMissing
	}
	server, err := s.servers.Get(ctx, serverID)
	if err != nil {
		zap.L().Error("failed to get server for deployment",
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		return "", ErrApplicationMissing
	}
	deployments, err := s.repository.List(ctx, serverID)
	if err != nil {
		zap.L().Error("failed to list existing deployments for conflict check",
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	for _, deployment := range deployments {
		existing, err := s.applications.Get(ctx, deployment.ApplicationID)
		if err != nil {
			zap.L().Warn("failed to get existing application for conflict check",
				zap.Uint("application_id", deployment.ApplicationID),
				zap.Uint("server_id", serverID),
				zap.Uint64("deploy_id", deployment.DeployID),
				zap.Error(err),
			)
			continue
		}
		if existing.Type == "compose" && app.Type == "compose" {
			if err := checkComposeContent(app.Content, deployment.Content); err != nil {
				zap.L().Warn("deployment compose check failed",
					zap.Uint("application_id", applicationID),
					zap.String("application_name", app.Name),
					zap.Uint("existing_application_id", deployment.ApplicationID),
					zap.String("existing_application_name", existing.Name),
					zap.Uint("server_id", serverID),
					zap.Uint64("existing_deploy_id", deployment.DeployID),
					zap.Error(err),
				)
				return "", err
			}
		}
	}
	deployID, err := s.ids.Generate()
	if err != nil {
		zap.L().Error("failed to generate deployment ID",
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		return "", ErrIDGeneration
	}
	request := AgentApplication{
		Name: app.Name, Description: app.Description, Type: app.Type, Content: app.Content,
		Version: app.Version, ServerID: serverID, DeployID: deployID,
	}
	if err := s.agent.Post(ctx, server, "application", request); err != nil {
		zap.L().Error("failed to deploy application on agent",
			zap.Uint64("deploy_id", deployID),
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		return "", ErrAgentDeploy
	}
	deployment := domain.Deployment{
		ServerID: serverID, ApplicationID: applicationID, Content: app.Content, DeployID: deployID,
	}
	if err := s.repository.Add(ctx, &deployment); err != nil {
		zap.L().Error("failed to create deployment record",
			zap.Uint64("deploy_id", deployID),
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
			zap.Error(err),
		)
		// Preserve the old rollback path bug: it uses the application name,
		// while normal undeploy uses the numeric deploy ID.
		rollbackPath := fmt.Sprintf("application/delete/%s", app.Name)
		zap.L().Info("attempting to roll back application deployment",
			zap.Uint64("deploy_id", deployID),
			zap.Uint("application_id", applicationID),
			zap.Uint("server_id", serverID),
		)
		if rollbackErr := s.agent.Post(ctx, server, rollbackPath, nil); rollbackErr != nil {
			zap.L().Error("failed to roll back application deployment",
				zap.Uint64("deploy_id", deployID),
				zap.Uint("application_id", applicationID),
				zap.Uint("server_id", serverID),
				zap.Error(rollbackErr),
			)
		}
		return "", repositoryError(err)
	}
	return "deploy success", nil
}

func (s *Service) ReDeploy(ctx context.Context, id uint) (string, error) {
	deployment, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get deployment for redeploy", zap.Uint("deployment_id", id), zap.Error(err))
		return "", repositoryError(err)
	}
	app, err := s.applications.Get(ctx, deployment.ApplicationID)
	if err != nil {
		zap.L().Error("failed to get application for redeploy",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("application_id", deployment.ApplicationID),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	server, err := s.servers.Get(ctx, deployment.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for redeploy",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return "", ErrApplicationMissing
	}
	request := AgentApplication{
		Name: app.Name, Description: app.Description, Type: app.Type, Content: deployment.Content,
		Version: app.Version, ServerID: deployment.ServerID, DeployID: deployment.DeployID,
	}
	if err := s.agent.Post(ctx, server, "application", request); err != nil {
		zap.L().Error("failed to redeploy application on agent",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("application_id", deployment.ApplicationID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return "", ErrAgentDeploy
	}
	return "success", nil
}

func (s *Service) List(ctx context.Context, serverID uint) ([]DeploymentView, error) {
	deployments, err := s.repository.List(ctx, serverID)
	if err != nil {
		zap.L().Error("failed to list deployments", zap.Uint("server_id", serverID), zap.Error(err))
		return nil, repositoryError(err)
	}
	var result []DeploymentView
	for _, deployment := range deployments {
		app, err := s.applications.Get(ctx, deployment.ApplicationID)
		if err != nil {
			zap.L().Warn("failed to get application information for deployment",
				zap.Uint("application_id", deployment.ApplicationID),
				zap.Uint("server_id", deployment.ServerID),
				zap.Uint64("deploy_id", deployment.DeployID),
				zap.Error(err),
			)
			continue
		}
		server, err := s.servers.Get(ctx, deployment.ServerID)
		if err != nil {
			zap.L().Warn("failed to get server information for deployment",
				zap.Uint("application_id", deployment.ApplicationID),
				zap.Uint("server_id", deployment.ServerID),
				zap.Uint64("deploy_id", deployment.DeployID),
				zap.Error(err),
			)
			continue
		}
		result = append(result, DeploymentView{Deployment: deployment, Application: app, Server: server})
	}
	return result, nil
}

func (s *Service) ListServers(ctx context.Context, applicationID uint) ([]domain.Server, error) {
	if _, err := s.applications.Get(ctx, applicationID); err != nil {
		zap.L().Error("failed to get application for listing servers",
			zap.Uint("application_id", applicationID),
			zap.Error(err),
		)
		return nil, ErrApplicationMissing
	}
	// Preserve the old repository call: List interprets this argument as
	// server_id, despite this endpoint passing applicationID.
	deployments, err := s.repository.List(ctx, applicationID)
	if err != nil {
		zap.L().Error("failed to list application servers",
			zap.Uint("application_id", applicationID),
			zap.Error(err),
		)
		return nil, repositoryError(err)
	}
	var result []domain.Server
	for _, deployment := range deployments {
		server, err := s.servers.Get(ctx, deployment.ServerID)
		if err != nil {
			zap.L().Warn("failed to get server information",
				zap.Uint("application_id", applicationID),
				zap.Uint("server_id", deployment.ServerID),
				zap.Uint64("deploy_id", deployment.DeployID),
				zap.Error(err),
			)
			continue
		}
		result = append(result, server)
	}
	return result, nil
}

func (s *Service) Update(ctx context.Context, id uint, content string) (string, error) {
	deployment, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get deployment for update", zap.Uint("deployment_id", id), zap.Error(err))
		return "", repositoryError(err)
	}
	deployment.Content = content
	if err := s.repository.Update(ctx, &deployment); err != nil {
		zap.L().Error("failed to update deployment",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	return "success", nil
}

func (s *Service) Undeploy(ctx context.Context, id uint) (string, error) {
	deployment, server, err := s.deploymentServer(ctx, id)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("application/delete/%d", deployment.DeployID)
	if err := s.agent.Post(ctx, server, path, nil); err != nil {
		zap.L().Error("failed to delete application on agent",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return "", ErrAgentDelete
	}
	if err := s.repository.Delete(ctx, deployment.ID); err != nil {
		zap.L().Error("failed to delete deployment record",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	return "undeploy success", nil
}

func (s *Service) Stop(ctx context.Context, id uint) (string, error) {
	deployment, server, err := s.deploymentServer(ctx, id)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("application/stop/%d", deployment.DeployID)
	if err := s.agent.Post(ctx, server, path, nil); err != nil {
		zap.L().Error("failed to stop application on agent",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return "", ErrAgentStop
	}
	return "stop success", nil
}

func (s *Service) Start(ctx context.Context, id uint) (string, error) {
	deployment, server, err := s.deploymentServer(ctx, id)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("application/start/%d", deployment.DeployID)
	if err := s.agent.Post(ctx, server, path, nil); err != nil {
		zap.L().Error("failed to start application on agent",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return "", ErrAgentStart
	}
	return "start success", nil
}

func (s *Service) ReportStatus(ctx context.Context, deployID uint64, status string) (string, error) {
	if _, err := s.repository.GetByDeployID(ctx, deployID); err != nil {
		zap.L().Error("failed to get deployment for status update",
			zap.Uint64("deploy_id", deployID),
			zap.String("status", status),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	if err := s.repository.UpdateStatus(ctx, deployID, status); err != nil {
		zap.L().Error("failed to update application status",
			zap.Uint64("deploy_id", deployID),
			zap.String("status", status),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	zap.L().Info("application status updated", zap.Uint64("deploy_id", deployID), zap.String("status", status))
	return "success", nil
}

func (s *Service) deploymentServer(ctx context.Context, id uint) (domain.Deployment, domain.Server, error) {
	deployment, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get deployment", zap.Uint("deployment_id", id), zap.Error(err))
		return domain.Deployment{}, domain.Server{}, repositoryError(err)
	}
	server, err := s.servers.Get(ctx, deployment.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for deployment operation",
			zap.Uint("deployment_id", id),
			zap.Uint64("deploy_id", deployment.DeployID),
			zap.Uint("server_id", deployment.ServerID),
			zap.Error(err),
		)
		return domain.Deployment{}, domain.Server{}, ErrApplicationMissing
	}
	return deployment, server, nil
}
