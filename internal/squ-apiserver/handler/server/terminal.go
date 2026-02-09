package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
	"squirrel-dev/internal/squ-apiserver/terminal"
	"squirrel-dev/pkg/jwt"
	"squirrel-dev/pkg/ssh"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// AuthMessage 认证消息结构
type AuthMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (s *Server) GetTerminal(id uint, conn *websocket.Conn) response.Response {
	// 1. 先等待认证消息
	var authMsg AuthMessage
	if err := conn.ReadJSON(&authMsg); err != nil {
		zap.L().Error("failed to read auth message", zap.Error(err))
		terminal.WriteMessage(conn, "error", "failed to read auth message")
		conn.Close()
		return response.Error(res.ErrConnectFailed)
	}

	// 验证消息类型
	if authMsg.Type != "auth" {
		zap.L().Warn("invalid message type, expected auth", zap.String("type", authMsg.Type))
		terminal.WriteMessage(conn, "error", "expected auth message")
		conn.Close()
		return response.Error(res.ErrConnectFailed)
	}

	// 验证 token
	j := jwt.New(s.Config.Auth.Jwt.SigningKey)
	claims, err := j.ParseToken(authMsg.Token)
	if err != nil {
		zap.L().Warn("invalid token", zap.Error(err))
		terminal.WriteMessage(conn, "auth_failed", "invalid token")
		conn.Close()
		return response.Error(res.ErrConnectFailed)
	}

	// 发送认证成功消息
	if err := terminal.WriteMessage(conn, "auth_success", "authenticated"); err != nil {
		zap.L().Error("failed to send auth success", zap.Error(err))
		conn.Close()
		return response.Error(res.ErrConnectFailed)
	}

	zap.L().Info("terminal websocket authenticated", zap.String("username", claims.Username), zap.Uint("server_id", id))

	// 2. 认证通过，获取服务器信息并建立 SSH 连接
	daoS, err := s.Repository.Get(id)
	if err != nil {
		terminal.WriteMessage(conn, "error", "server not found")
		conn.Close()
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
		zap.L().Error("failed to establish ssh connection",
			zap.String("ip_address", daoS.IpAddress),
			zap.String("username", daoS.SshUsername),
			zap.Error(err),
		)
		terminal.WriteMessage(conn, "error", "failed to connect to server")
		conn.Close()
		return response.Error(res.ErrConnectFailed)
	}

	terminalHandler, err := terminal.NewTerminalHandler("ssh", 80, 24, sshClient.Client)
	if err != nil {
		zap.L().Error("failed to initialize terminal", zap.Error(err))
		terminal.WriteMessage(conn, "error", "failed to initialize terminal")
		conn.Close()
		return response.Error(res.ErrConnectFailed)
	}

	// 3. 启动 WebSocket 终端处理
	terminal.HandleWebSocket(conn, terminalHandler)
	return response.Success("success")
}
