package health

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	healthRepository "squirrel-dev/internal/squ-apiserver/repository/health"
)

type Health struct {
	Config     *config.Config
	Repository healthRepository.Client
}

func (h *Health) Status() response.Response {
	return response.Success("health")
}
