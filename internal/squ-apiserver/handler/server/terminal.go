package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/ssh"
	"squirrel-dev/pkg/terminal"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *Server) GetTerminal(id uint, conn *websocket.Conn) response.Response {
	daoS, err := s.ModelClient.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
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
	err = terminal.HandleWebSocketWithSSH(conn, sshClient.Client)
	if err != nil {
		zap.S().Error("ssh connect failed", zap.Error(err))
		return response.Error(res.ErrConnectFailed)
	}
	return response.Success("success")
}
