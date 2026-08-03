package api

import (
	"squirrel-dev/internal/squ-apiserver/module/deployment/api/res"
	"squirrel-dev/internal/squ-apiserver/module/deployment/application"
	"squirrel-dev/internal/squ-apiserver/module/deployment/domain"
)

func toDeploymentResponse(value application.DeploymentView) res.Deployment {
	return res.Deployment{
		ID:       value.Deployment.ID,
		DeployID: value.Deployment.DeployID,
		Application: res.ApplicationInfo{
			ID:          value.Application.ID,
			Name:        value.Application.Name,
			Description: value.Application.Description,
			Type:        value.Application.Type,
			Version:     value.Application.Version,
		},
		Server:     toServerResponse(value.Server),
		Status:     value.Deployment.Status,
		DeployedAt: value.Deployment.CreatedAt.Format("2006-01-02 15:04:05"),
		Content:    value.Deployment.Content,
	}
}

func toServerResponse(value domain.Server) res.ServerInfo {
	return res.ServerInfo{
		ID:        value.ID,
		IPAddress: value.IPAddress,
		AgentPort: value.AgentPort,
	}
}
