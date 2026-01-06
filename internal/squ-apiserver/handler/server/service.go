package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/model"

	serverModel "squirrel-dev/internal/squ-apiserver/model/server"
)

type Server struct {
	Config      *config.Config
	ModelClient serverModel.Client
}

func (s *Server) List() response.Response {
	var servers []res.Server
	daoServers, err := s.ModelClient.List()
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
			Status:      daoS.Status,
		})
	}
	return response.Success(servers)
}

func (s *Server) Get(id uint) response.Response {
	var serverRes res.Server
	daoS, err := s.ModelClient.Get(id)
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

	return response.Success(serverRes)
}

func (s *Server) Delete(id uint) response.Response {
	err := s.ModelClient.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Add(request req.Server) response.Response {
	modelReq := serverModel.Server{
		Hostname:    request.Hostname,
		IpAddress:   request.IpAddress,
		SshUsername: request.SshUsername,
		SshPort:     request.SshPort,
		AuthType:    request.AuthType,
		Status:      request.Status,
	}

	err := s.ModelClient.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Update(request req.Server) response.Response {
	modelReq := serverModel.Server{
		IpAddress:   request.IpAddress,
		SshUsername: request.SshUsername,
		SshPort:     request.SshPort,
		AuthType:    request.AuthType,
		Status:      request.Status,
	}
	modelReq.ID = request.ID
	err := s.ModelClient.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
