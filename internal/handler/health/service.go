package health

import (
	"squirrel-dev/internal/config"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
	healthModel "squirrel-dev/internal/model/health"
)

type Health struct{
	Config *config.Config
	DB     database.DB
	HealthModelClient healthModel.Client
}

func (h *Health) Status() response.Response {
	return response.Success("health")
}
