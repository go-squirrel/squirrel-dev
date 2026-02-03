package server

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"
	"time"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
)

type Server struct {
	Config     *config.Config
	Repository serverRepository.Repository
	HTTPClient *httpclient.Client
}

func New(conf *config.Config, repo serverRepository.Repository) *Server {
	hc := httpclient.NewClient(30 * time.Second)
	return &Server{
		Config:     conf,
		Repository: repo,
		HTTPClient: hc,
	}
}

func (s *Server) List() response.Response {
	var servers []res.Server
	daoServers, err := s.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoS := range daoServers {
		servers = append(servers, res.Server{
			ID:          daoS.ID,
			Hostname:    daoS.Hostname,
			IpAddress:   daoS.IpAddress,
			SshUsername: daoS.SshUsername,
			SshPort:     daoS.SshPort,
			AuthType:    daoS.AuthType,
			Status:      model.ServerStatusOnline,
		})
	}
	return response.Success(servers)
}

func (s *Server) Get(id uint) response.Response {
	var serverRes res.Server
	daoS, err := s.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	serverRes = res.Server{
		ID:          daoS.ID,
		Hostname:    daoS.Hostname,
		IpAddress:   daoS.IpAddress,
		SshUsername: daoS.SshUsername,
		SshPort:     daoS.SshPort,
		AuthType:    daoS.AuthType,
		Status:      daoS.Status,
	}

	agentURL := utils.GenAgentUrl(s.Config.Agent.Http.Scheme,
		daoS.IpAddress,
		daoS.AgentPort,
		s.Config.Agent.Http.BaseUrl,
		"server/info")

	respBody, err := s.HTTPClient.Get(agentURL, nil)
	var agentResp response.Response
	if err != nil {
		zap.L().Error("获取 Agent 信息失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		serverRes.Status = model.ServerStatusOffline
	} else {
		if err := json.Unmarshal(respBody, &agentResp); err != nil {
			zap.L().Error("解析 Agent 响应失败",
				zap.String("url", agentURL),
				zap.Error(err),
			)
			serverRes.Status = model.ServerStatusOffline
		}
		if agentResp.Code != 0 {
			zap.L().Error("Agent 获取信息失败",
				zap.String("url", agentURL),
				zap.Int("code", agentResp.Code),
				zap.String("message", agentResp.Message),
			)
			serverRes.Status = model.ServerStatusOffline
		}
	}
	if agentResp.Data != nil {
		serverRes.ServerInfo = agentResp.Data.(map[string]any)
	}

	return response.Success(serverRes)
}

func (s *Server) Delete(id uint) response.Response {
	err := s.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Add(request req.Server) response.Response {
	modelReq := model.Server{
		Hostname:    request.Hostname,
		IpAddress:   request.IpAddress,
		SshUsername: request.SshUsername,
		SshPort:     request.SshPort,
		AuthType:    request.AuthType,
		Status:      request.Status,
	}

	err := s.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Update(request req.Server) response.Response {
	modelReq := model.Server{
		IpAddress:   request.IpAddress,
		SshUsername: request.SshUsername,
		SshPort:     request.SshPort,
		AuthType:    request.AuthType,
		Status:      request.Status,
	}
	modelReq.ID = request.ID
	err := s.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
