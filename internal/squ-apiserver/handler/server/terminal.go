package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/pkg/ssh"
	"squirrel-dev/pkg/terminal"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *Server) GetTerminal(id uint, conn *websocket.Conn) response.Response {
	daoS, err := s.Repository.Get(id)
	if err != nil {
		return response.Error(returnServerErrCode(err))
	}
	privateKey := ""
	if daoS.SshPrivateKey != nil {
		privateKey = *daoS.SshPrivateKey
	}
	password := ""
	if daoS.SshPassword != nil {
		password = *daoS.SshPassword
	}
	machine := &ssh.Machine{
		Name:       "test",
		IpAddress:  daoS.IpAddress,
		User:       daoS.SshUsername,
		Password:   password,
		Port:       daoS.SshPort,
		PrivateKey: privateKey,
		Type:       daoS.AuthType,
	}
	sshClient, err := ssh.NewSsh(machine)
	if err != nil {
		zap.S().Error("ssh connect failed", zap.Error(err))
		return response.Error(res.ErrConnectFailed)
	}
	terminalHandler, err := terminal.NewTerminalHandler("ssh", 80, 24, sshClient.Client)
	if err != nil {
		zap.S().Error("terminal init failed", zap.Error(err))
		return response.Error(res.ErrConnectFailed)
	}
	terminal.HandleWebSocket(conn, terminalHandler)
	return response.Success("success")
}
