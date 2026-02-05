package server

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

func (s *Server) getAgentInfo(ipAddress string, port int) (status string, agentResp response.Response) {
	agentURL := utils.GenAgentUrl(s.Config.Agent.Http.Scheme,
		ipAddress,
		port,
		s.Config.Agent.Http.BaseUrl,
		"server/info")

	respBody, err := s.HTTPClient.Get(agentURL, nil)

	if err != nil {
		zap.L().Error("获取 Agent 信息失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return model.ServerStatusOffline, agentResp
	} else {
		if err := json.Unmarshal(respBody, &agentResp); err != nil {
			zap.L().Error("解析 Agent 响应失败",
				zap.String("url", agentURL),
				zap.Error(err),
			)
			return model.ServerStatusOffline, agentResp
		}
		if agentResp.Code != 0 {
			zap.L().Error("Agent 获取信息失败",
				zap.String("url", agentURL),
				zap.Int("code", agentResp.Code),
				zap.String("message", agentResp.Message),
			)
			return model.ServerStatusOffline, agentResp
		}
	}
	return model.ServerStatusOnline, agentResp
}
