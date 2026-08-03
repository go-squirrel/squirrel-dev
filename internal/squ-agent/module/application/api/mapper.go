package api

import (
	"squirrel-dev/internal/squ-agent/module/application/api/req"
	"squirrel-dev/internal/squ-agent/module/application/api/res"
	"squirrel-dev/internal/squ-agent/module/application/application"
	"squirrel-dev/internal/squ-agent/module/application/domain"
)

func toApplicationRequest(value req.Application) application.Request {
	return application.Request{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Type:        value.Type,
		Status:      value.Status,
		Content:     value.Content,
		Version:     value.Version,
		ServerID:    value.ServerID,
		DeployID:    value.DeployID,
		Env:         value.Env,
	}
}

func fromDomain(value domain.Application) res.Application {
	return res.Application{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Type:        value.Type,
		Status:      value.Status,
		Content:     value.Content,
		Version:     value.Version,
	}
}
