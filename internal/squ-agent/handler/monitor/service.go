package monitor

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	monitorRepository "squirrel-dev/internal/squ-agent/repository/monitor"
)

type Monitor struct {
	Config     *config.Config
	Repository monitorRepository.Client
}

func (m *Monitor) Status() response.Response {
	return response.Success("monitor")
}
