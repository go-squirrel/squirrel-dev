package monitor

import (
	"encoding/json"
	"fmt"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/monitor/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
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
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		"monitor/stats")

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取监控数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetDiskIO(serverID uint, device string) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		"monitor/stats/io/"+device)

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取磁盘IO数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetAllDiskIO(serverID uint) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		"monitor/stats/io/all")

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取磁盘IO数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetNetIO(serverID uint, interfaceName string) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		"monitor/stats/net/"+interfaceName)

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取网络IO数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetAllNetIO(serverID uint) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		"monitor/stats/net/all")

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取网络IO数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetBaseMonitorPage(serverID uint, page, count int) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	baseMonitorUrl := fmt.Sprintf("monitor/base/%d/%d", page, count)
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		baseMonitorUrl)
	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取基础监控数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetDiskIOMonitorPage(serverID uint, page, count int) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	diskMonitorUrl := fmt.Sprintf("monitor/disk/%d/%d", page, count)
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		diskMonitorUrl)

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取磁盘监控数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}

func (m *Monitor) GetNetworkMonitorPage(serverID uint, page, count int) (response.Response, error) {
	// 检查服务器是否存在
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound), nil
		}
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 构建请求URL
	netMonitorUrl := fmt.Sprintf("monitor/net/%d/%d", page, count)
	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		netMonitorUrl)

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用 Agent 获取网络监控数据失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}
