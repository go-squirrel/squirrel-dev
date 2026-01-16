package health

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	healthRepository "squirrel-dev/internal/squ-agent/repository/health"
)

type Health struct {
	Config     *config.Config
	DB         database.DB
	Repository healthRepository.Client
}

func (h *Health) Status() response.Response {
	return response.Success("health")
}
