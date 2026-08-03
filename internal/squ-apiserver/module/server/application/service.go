package application

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/server/domain"
)

type Request struct {
	ID            uint
	Hostname      string
	IPAddress     string
	Port          int
	SSHUsername   string
	SSHPassword   string
	SSHPrivateKey string
	SSHPort       int
	AuthType      string
	Status        string
	ServerAlias   string
}

type ServerView struct {
	ServerInfo map[string]any
	domain.Server
}

type Service struct {
	repository domain.Repository
	agents     domain.AgentInfoClient
	ssh        domain.SSHTester
}

func NewService(repository domain.Repository, agents domain.AgentInfoClient, ssh domain.SSHTester) *Service {
	return &Service{repository: repository, agents: agents, ssh: ssh}
}

func (s *Service) List(ctx context.Context) ([]ServerView, error) {
	servers, err := s.repository.List(ctx)
	if err != nil {
		zap.L().Error("failed to list servers", zap.Error(err))
		return nil, err
	}
	var result []ServerView
	for _, server := range servers {
		status, _ := s.agents.GetInfo(ctx, server.IPAddress, server.AgentPort)
		server.Status = status
		result = append(result, ServerView{Server: server})
	}
	return result, nil
}

func (s *Service) Get(ctx context.Context, id uint) (ServerView, error) {
	server, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get server", zap.Uint("server_id", id), zap.Error(err))
		return ServerView{}, err
	}
	status, info := s.agents.GetInfo(ctx, server.IPAddress, server.AgentPort)
	server.Status = status
	return ServerView{Server: server, ServerInfo: info}, nil
}

func (s *Service) GetStored(ctx context.Context, id uint) (domain.Server, error) {
	server, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get stored server", zap.Uint("server_id", id), zap.Error(err))
		return domain.Server{}, err
	}
	return server, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		zap.L().Error("failed to delete server", zap.Uint("server_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Add(ctx context.Context, request Request) error {
	server := requestToServer(request)
	server.UUID = uuid.New().String()
	if err := s.repository.Add(ctx, &server); err != nil {
		zap.L().Error("failed to add server",
			zap.String("hostname", server.Hostname),
			zap.String("ip_address", server.IPAddress),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, request Request) error {
	existing, err := s.repository.Get(ctx, request.ID)
	if err != nil {
		zap.L().Error("failed to get server for update", zap.Uint("server_id", request.ID), zap.Error(err))
		return err
	}
	server := requestToServer(request)
	server.ID = request.ID
	server.UUID = existing.UUID
	if err := s.repository.Update(ctx, &server); err != nil {
		zap.L().Error("failed to update server",
			zap.Uint("server_id", request.ID),
			zap.String("hostname", server.Hostname),
			zap.String("ip_address", server.IPAddress),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *Service) CheckAgent(ctx context.Context, ip string, port int) (bool, string, map[string]any) {
	status, info := s.agents.GetInfo(ctx, ip, port)
	if status == domain.StatusOnline {
		return true, "Agent is ready", info
	}
	return false, "Agent is not ready", nil
}

func (s *Service) TestSSH(ctx context.Context, id uint) (domain.Server, error) {
	server, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get server for ssh test", zap.Uint("server_id", id), zap.Error(err))
		return domain.Server{}, err
	}
	if err := s.ssh.Test(ctx, server); err != nil {
		zap.L().Error("ssh connection test failed",
			zap.Uint("server_id", id),
			zap.String("ip_address", server.IPAddress),
			zap.String("username", server.SSHUsername),
			zap.Error(err),
		)
		return domain.Server{}, err
	}
	zap.L().Info("ssh connection test successful",
		zap.Uint("server_id", id),
		zap.String("ip_address", server.IPAddress),
	)
	return server, nil
}

func requestToServer(request Request) domain.Server {
	hostname := request.Hostname
	if hostname == "" {
		hostname = request.IPAddress
	}
	server := domain.Server{
		Hostname: request.Hostname, IPAddress: request.IPAddress, AgentPort: request.Port,
		SSHUsername: request.SSHUsername, SSHPort: request.SSHPort,
		AuthType: request.AuthType, Status: request.Status,
	}
	server.Hostname = hostname
	if request.ServerAlias != "" {
		server.ServerAlias = &request.ServerAlias
	}
	if request.AuthType == domain.AuthTypePassword {
		if request.SSHPassword != "" {
			server.SSHPassword = &request.SSHPassword
		}
	} else if request.SSHPrivateKey != "" {
		server.SSHPrivateKey = &request.SSHPrivateKey
	}
	return server
}
