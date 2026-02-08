package server

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
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

func (s *Server) modelToResponse(daoS model.Server, status string) res.Server {
	return res.Server{
		ID:            daoS.ID,
		Hostname:      daoS.Hostname,
		IpAddress:     daoS.IpAddress,
		Port:          daoS.AgentPort,
		SshUsername:   daoS.SshUsername,
		SshPassword:   daoS.SshPassword,
		SshPrivateKey: daoS.SshPrivateKey,
		ServerAlias:   daoS.ServerAlias,
		SshPort:       daoS.SshPort,
		AuthType:      daoS.AuthType,
		Status:        status,
	}
}

func (s *Server) requestToModel(request req.Server) model.Server {
	modelReq := model.Server{
		Hostname:    request.Hostname,
		IpAddress:   request.IpAddress,
		SshUsername: request.SshUsername,
		SshPort:     request.SshPort,
		AuthType:    request.AuthType,
		Status:      request.Status,
	}
	if request.AuthType == model.ServerAuthTypePassword {
		modelReq.SshPassword = &request.SshPassword
	} else {
		modelReq.SshPrivateKey = &request.SshPrivateKey
	}
	return modelReq
}

// returnServerErrCode 根据错误类型返回精确的服务器错误码
func returnServerErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrServerNotFound
	case gorm.ErrDuplicatedKey:
		return res.ErrServerAlreadyExists
	}
	return res.ErrServerUpdateFailed
}
