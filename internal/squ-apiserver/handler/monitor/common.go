package monitor

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/monitor/res"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (m *Monitor) callAgent(serverID uint, path string, description string) response.Response {
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		return response.Error(returnMonitorErrCode(err))
	}

	agentURL := utils.GenAgentUrl(m.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		m.Config.Agent.Http.BaseUrl,
		path)

	respBody, err := m.HTTPClient.Get(agentURL, nil)
	if err != nil {
		zap.L().Error("调用Agent失败",
			zap.String("description", description),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed)
	}

	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析Agent响应失败",
			zap.String("description", description),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed)
	}

	return agentResp
}

// returnMonitorErrCode 根据错误类型返回精确的监控错误码
func returnMonitorErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrServerNotFound
	}
	return res.ErrMonitorFailed
}
