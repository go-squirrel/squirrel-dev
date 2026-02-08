package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/httpclient"
	"time"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
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

		status, _ := s.getAgentInfo(daoS.IpAddress, daoS.AgentPort)

		servers = append(servers, s.modelToResponse(daoS, status))
	}
	return response.Success(servers)
}

func (s *Server) Get(id uint) response.Response {
	daoS, err := s.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	status, agentResp := s.getAgentInfo(daoS.IpAddress, daoS.AgentPort)

	serverRes := s.modelToResponse(daoS, status)
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
	modelReq := s.requestToModel(request)

	err := s.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Update(request req.Server) response.Response {
	modelReq := s.requestToModel(request)
	modelReq.ID = request.ID

	err := s.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Registry(request req.Register) response.Response {
	daoS, err := s.Repository.GetByUUID(request.UUID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	daoS.AgentPort = request.AgentPort
	err = s.Repository.Update(&daoS)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
