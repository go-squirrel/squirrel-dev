package monitor

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/monitor/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

func (m *Monitor) callAgent(serverID uint, path string, description string) (response.Response, error) {
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err)), nil
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
		return response.Error(res.ErrMonitorFailed), nil
	}

	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析Agent响应失败",
			zap.String("description", description),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return response.Error(res.ErrMonitorFailed), nil
	}

	return agentResp, nil
}
