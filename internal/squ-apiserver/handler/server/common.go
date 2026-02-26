package server

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"github.com/google/uuid"
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
		zap.L().Error("failed to get agent information",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		return model.ServerStatusOffline, agentResp
	} else {
		if err := json.Unmarshal(respBody, &agentResp); err != nil {
			zap.L().Error("failed to parse agent response",
				zap.String("url", agentURL),
				zap.Error(err),
			)
			return model.ServerStatusOffline, agentResp
		}
		if agentResp.Code != 0 {
			zap.L().Error("agent failed to get information",
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
	// 如果 hostname 为空，使用 ip_address 作为默认值
	hostname := request.Hostname
	if hostname == "" {
		hostname = request.IpAddress
	}

	modelReq := model.Server{
		Hostname:    hostname,
		IpAddress:   request.IpAddress,
		AgentPort:   request.Port,
		SshUsername: request.SshUsername,
		SshPort:     request.SshPort,
		AuthType:    request.AuthType,
		Status:      request.Status,
	}

	// 处理 server_alias
	if request.ServerAlias != "" {
		modelReq.ServerAlias = &request.ServerAlias
	}

	// 处理认证信息
	if request.AuthType == model.ServerAuthTypePassword {
		if request.SshPassword != "" {
			modelReq.SshPassword = &request.SshPassword
		}
	} else {
		if request.SshPrivateKey != "" {
			modelReq.SshPrivateKey = &request.SshPrivateKey
		}
	}
	return modelReq
}

// generateUUID 生成服务器唯一标识
func generateUUID() string {
	return uuid.New().String()
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
