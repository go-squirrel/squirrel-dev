package monitor

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	monitorModel "squirrel-dev/internal/squ-agent/model/monitor"
)

type Monitor struct {
	Config             *config.Config
	DB                 database.DB
	MonitorModelClient monitorModel.Client
}

func (m *Monitor) Status() response.Response {
	return response.Success("monitor")
}
