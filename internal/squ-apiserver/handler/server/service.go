package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"

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
		return response.Error(response.ErrSQLNotFound)
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
