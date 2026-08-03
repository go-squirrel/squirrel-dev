package api

import (
	"squirrel-dev/internal/squ-apiserver/module/server/api/req"
	"squirrel-dev/internal/squ-apiserver/module/server/api/res"
	"squirrel-dev/internal/squ-apiserver/module/server/application"
	"squirrel-dev/internal/squ-apiserver/module/server/domain"
)

func toApplication(value req.Server) application.Request {
	return application.Request{
		ID:            value.ID,
		Hostname:      value.Hostname,
		IPAddress:     value.IPAddress,
		Port:          value.Port,
		SSHUsername:   value.SSHUsername,
		SSHPassword:   value.SSHPassword,
		SSHPrivateKey: value.SSHPrivateKey,
		SSHPort:       value.SSHPort,
		AuthType:      value.AuthType,
		Status:        value.Status,
		ServerAlias:   value.ServerAlias,
	}
}

func toResponse(value application.ServerView) res.Server {
	server := value.Server
	return res.Server{
		ID:            server.ID,
		Hostname:      server.Hostname,
		IPAddress:     server.IPAddress,
		Port:          server.AgentPort,
		SSHUsername:   server.SSHUsername,
		SSHPassword:   server.SSHPassword,
		SSHPrivateKey: server.SSHPrivateKey,
		SSHPort:       server.SSHPort,
		AuthType:      server.AuthType,
		Status:        server.Status,
		ServerAlias:   server.ServerAlias,
		ServerInfo:    value.ServerInfo,
	}
}

func toSSHTestResponse(value domain.Server) res.SSHTestResult {
	return res.SSHTestResult{
		Message:   "SSH connection successful",
		Hostname:  value.Hostname,
		IPAddress: value.IPAddress,
		SSHPort:   value.SSHPort,
	}
}

func toAgentCheckResponse(ready bool, message string, serverInfo map[string]any) res.AgentCheckResult {
	return res.AgentCheckResult{
		Ready:      ready,
		Message:    message,
		ServerInfo: serverInfo,
	}
}
