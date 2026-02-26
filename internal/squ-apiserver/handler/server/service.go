package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/ssh"
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
	hc := httpclient.NewClient(3 * time.Second)
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
		zap.L().Error("failed to list servers",
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
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
		zap.L().Error("failed to get server",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
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
		zap.L().Error("failed to delete server",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Add(request req.Server) response.Response {
	modelReq := s.requestToModel(request)

	err := s.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("failed to add server",
			zap.String("hostname", request.Hostname),
			zap.String("ip_address", request.IpAddress),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Update(request req.Server) response.Response {
	modelReq := s.requestToModel(request)
	modelReq.ID = request.ID

	err := s.Repository.Update(&modelReq)

	if err != nil {
		zap.L().Error("failed to update server",
			zap.Uint("id", request.ID),
			zap.String("hostname", request.Hostname),
			zap.String("ip_address", request.IpAddress),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
	}

	return response.Success("success")
}

func (s *Server) Registry(request req.Register) response.Response {
	daoS, err := s.Repository.GetByUUID(request.UUID)
	if err != nil {
		zap.L().Error("failed to get server by UUID",
			zap.String("uuid", request.UUID),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
	}

	daoS.AgentPort = request.AgentPort
	err = s.Repository.Update(&daoS)
	if err != nil {
		zap.L().Error("failed to update server agent port",
			zap.Uint("id", daoS.ID),
			zap.String("uuid", request.UUID),
			zap.Int("agent_port", request.AgentPort),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
	}

	return response.Success("success")
}

// TestSSH 测试 SSH 连接
func (s *Server) TestSSH(id uint) response.Response {
	// 从数据库获取服务器信息
	daoS, err := s.Repository.Get(id)
	if err != nil {
		zap.L().Error("failed to get server",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return response.Error(returnServerErrCode(err))
	}

	// 准备 SSH 连接参数
	privateKey := ""
	if daoS.SshPrivateKey != nil {
		privateKey = *daoS.SshPrivateKey
	}
	password := ""
	if daoS.SshPassword != nil {
		password = *daoS.SshPassword
	}

	machine := &ssh.Machine{
		Name:       daoS.Hostname,
		IpAddress:  daoS.IpAddress,
		User:       daoS.SshUsername,
		Password:   password,
		Port:       daoS.SshPort,
		PrivateKey: privateKey,
		Type:       daoS.AuthType,
	}

	// 尝试建立 SSH 连接
	sshClient, err := ssh.NewSsh(machine)
	if err != nil {
		zap.L().Error("SSH connection test failed",
			zap.Uint("id", id),
			zap.String("ip_address", daoS.IpAddress),
			zap.String("username", daoS.SshUsername),
			zap.Error(err),
		)
		return response.Error(res.ErrSSHTestFailed)
	}
	defer sshClient.Close()

	zap.L().Info("SSH connection test successful",
		zap.Uint("id", id),
		zap.String("ip_address", daoS.IpAddress),
	)

	return response.Success(map[string]any{
		"message":    "SSH connection successful",
		"hostname":   daoS.Hostname,
		"ip_address": daoS.IpAddress,
		"ssh_port":   daoS.SshPort,
	})
}
