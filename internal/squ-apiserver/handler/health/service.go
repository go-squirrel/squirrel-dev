package health

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	healthModel "squirrel-dev/internal/squ-apiserver/model/health"
)

type Health struct {
	Config            *config.Config
	DB                database.DB
	HealthModelClient healthModel.Client
}

func (h *Health) Status() response.Response {
	return response.Success("health")
}
