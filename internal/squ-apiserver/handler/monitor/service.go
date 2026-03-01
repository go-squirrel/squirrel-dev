package monitor

import (
	"fmt"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/agent"
	"squirrel-dev/internal/squ-apiserver/config"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

type Monitor struct {
	Config      *config.Config
	ServerRepo  serverRepository.Repository
	AgentClient *agent.Client
}

func New(conf *config.Config, serverRepo serverRepository.Repository) *Monitor {
	return &Monitor{
		Config:      conf,
		ServerRepo:  serverRepo,
		AgentClient: agent.NewClient(conf),
	}
}

func (m *Monitor) GetStats(serverID uint) response.Response {
	return m.callAgent(serverID, "monitor/stats", "get monitoring data")
}

func (m *Monitor) GetDiskIO(serverID uint, device string) response.Response {
	path := fmt.Sprintf("monitor/stats/io/%s", device)
	return m.callAgent(serverID, path, "get disk IO data")
}

func (m *Monitor) GetAllDiskIO(serverID uint) response.Response {
	return m.callAgent(serverID, "monitor/stats/io/all", "get all disk IO data")
}

func (m *Monitor) GetNetIO(serverID uint, interfaceName string) response.Response {
	path := fmt.Sprintf("monitor/stats/net/%s", interfaceName)
	return m.callAgent(serverID, path, "get network IO data")
}

func (m *Monitor) GetAllNetIO(serverID uint) response.Response {
	return m.callAgent(serverID, "monitor/stats/net/all", "get all network IO data")
}

func (m *Monitor) GetBaseMonitorByRange(serverID uint, timeRange string) response.Response {
	path := fmt.Sprintf("monitor/base?range=%s", timeRange)
	return m.callAgent(serverID, path, "get base monitor by time range")
}

func (m *Monitor) GetDiskIOMonitorByRange(serverID uint, timeRange string) response.Response {
	path := fmt.Sprintf("monitor/disk?range=%s", timeRange)
	return m.callAgent(serverID, path, "get disk IO monitor by time range")
}

func (m *Monitor) GetDiskUsageMonitorByRange(serverID uint, timeRange string) response.Response {
	path := fmt.Sprintf("monitor/disk-usage?range=%s", timeRange)
	return m.callAgent(serverID, path, "get disk usage monitor by time range")
}

func (m *Monitor) GetNetworkMonitorByRange(serverID uint, timeRange string) response.Response {
	path := fmt.Sprintf("monitor/net?range=%s", timeRange)
	return m.callAgent(serverID, path, "get network monitor by time range")
}
