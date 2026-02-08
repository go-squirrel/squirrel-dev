package monitor

import (
	"fmt"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/pkg/httpclient"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

type Monitor struct {
	Config     *config.Config
	ServerRepo serverRepository.Repository
	HTTPClient *httpclient.Client
}

func New(conf *config.Config, serverRepo serverRepository.Repository) *Monitor {
	hc := httpclient.NewClient(30 * time.Second)
	return &Monitor{
		Config:     conf,
		ServerRepo: serverRepo,
		HTTPClient: hc,
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

func (m *Monitor) GetBaseMonitorPage(serverID uint, page, count int) response.Response {
	path := fmt.Sprintf("monitor/base/%d/%d", page, count)
	return m.callAgent(serverID, path, "get base monitoring data page")
}

func (m *Monitor) GetDiskIOMonitorPage(serverID uint, page, count int) response.Response {
	path := fmt.Sprintf("monitor/disk/%d/%d", page, count)
	return m.callAgent(serverID, path, "get disk IO monitoring data page")
}

func (m *Monitor) GetNetworkMonitorPage(serverID uint, page, count int) response.Response {
	path := fmt.Sprintf("monitor/net/%d/%d", page, count)
	return m.callAgent(serverID, path, "get network monitoring data page")
}
