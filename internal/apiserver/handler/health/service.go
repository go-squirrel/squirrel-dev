package health

import (
	"squirrel-dev/internal/apiserver/config"
	healthModel "squirrel-dev/internal/apiserver/model/health"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
)

type Health struct {
	Config            *config.Config
	DB                database.DB
	HealthModelClient healthModel.Client
}

func (h *Health) Status() response.Response {
	return response.Success("health")
}
