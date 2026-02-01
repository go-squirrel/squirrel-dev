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

func (m *Monitor) GetStats(serverID uint) (response.Response, error) {
	return m.callAgent(serverID, "monitor/stats", "获取监控数据")
}

func (m *Monitor) GetDiskIO(serverID uint, device string) (response.Response, error) {
	path := fmt.Sprintf("monitor/stats/io/%s", device)
	return m.callAgent(serverID, path, "获取磁盘IO数据")
}

func (m *Monitor) GetAllDiskIO(serverID uint) (response.Response, error) {
	return m.callAgent(serverID, "monitor/stats/io/all", "获取所有磁盘IO数据")
}

func (m *Monitor) GetNetIO(serverID uint, interfaceName string) (response.Response, error) {
	path := fmt.Sprintf("monitor/stats/net/%s", interfaceName)
	return m.callAgent(serverID, path, "获取网络IO数据")
}

func (m *Monitor) GetAllNetIO(serverID uint) (response.Response, error) {
	return m.callAgent(serverID, "monitor/stats/net/all", "获取所有网络IO数据")
}

func (m *Monitor) GetBaseMonitorPage(serverID uint, page, count int) (response.Response, error) {
	path := fmt.Sprintf("monitor/base/%d/%d", page, count)
	return m.callAgent(serverID, path, "获取基础监控数据分页")
}

func (m *Monitor) GetDiskIOMonitorPage(serverID uint, page, count int) (response.Response, error) {
	path := fmt.Sprintf("monitor/disk/%d/%d", page, count)
	return m.callAgent(serverID, path, "获取磁盘IO监控数据分页")
}

func (m *Monitor) GetNetworkMonitorPage(serverID uint, page, count int) (response.Response, error) {
	path := fmt.Sprintf("monitor/net/%d/%d", page, count)
	return m.callAgent(serverID, path, "获取网络监控数据分页")
}
