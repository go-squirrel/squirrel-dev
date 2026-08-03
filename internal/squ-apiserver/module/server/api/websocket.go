package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/server/api/res"
	serverTerminal "squirrel-dev/internal/squ-apiserver/module/server/api/terminal"
	"squirrel-dev/internal/squ-apiserver/module/server/infra"
	"squirrel-dev/pkg/jwt"
	"squirrel-dev/pkg/utils"
)

type authMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (h *Handler) Terminal(c *gin.Context) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("failed to upgrade terminal websocket",
			zap.String("raw_server_id", c.Param("id")),
			zap.Error(err),
		)
		return
	}
	rawID := c.Param("id")
	id, err := parseServerID(rawID)
	if err != nil {
		zap.L().Warn("failed to parse terminal server ID", zap.String("raw_server_id", rawID), zap.Error(err))
		_ = conn.WriteJSON(response.Error(res.ErrInvalidParameter))
		_ = conn.Close()
		return
	}
	var auth authMessage
	if err := conn.ReadJSON(&auth); err != nil {
		zap.L().Error("failed to read terminal auth message", zap.Uint("server_id", id), zap.Error(err))
		_ = serverTerminal.WriteMessage(conn, "error", "failed to read auth message")
		_ = conn.Close()
		return
	}
	if auth.Type != "auth" {
		zap.L().Warn("invalid terminal websocket message type",
			zap.Uint("server_id", id),
			zap.String("type", auth.Type),
		)
		_ = serverTerminal.WriteMessage(conn, "error", "expected auth message")
		_ = conn.Close()
		return
	}
	claims, err := jwt.New(h.signingKey).ParseToken(auth.Token)
	if err != nil {
		zap.L().Warn("invalid terminal token", zap.Uint("server_id", id), zap.Error(err))
		_ = serverTerminal.WriteMessage(conn, "auth_failed", "invalid token")
		_ = conn.Close()
		return
	}
	if err := serverTerminal.WriteMessage(conn, "auth_success", "authenticated"); err != nil {
		zap.L().Error("failed to send terminal auth success", zap.Uint("server_id", id), zap.Error(err))
		_ = conn.Close()
		return
	}
	zap.L().Info("terminal websocket authenticated",
		zap.Uint("server_id", id),
		zap.String("username", claims.Username),
	)
	server, err := h.service.GetStored(c.Request.Context(), id)
	if err != nil {
		_ = serverTerminal.WriteMessage(conn, "error", "server not found")
		_ = conn.Close()
		return
	}
	client, err := infra.NewSSHClient(server)
	if err != nil {
		zap.L().Error("failed to establish terminal ssh connection",
			zap.Uint("server_id", id),
			zap.String("ip_address", server.IPAddress),
			zap.String("username", server.SSHUsername),
			zap.Error(err),
		)
		_ = serverTerminal.WriteMessage(conn, "error", "failed to connect to server")
		_ = conn.Close()
		return
	}
	terminalHandler, err := serverTerminal.NewSSH(client.Client, 80, 24)
	if err != nil {
		zap.L().Error("failed to initialize terminal", zap.Uint("server_id", id), zap.Error(err))
		_ = serverTerminal.WriteMessage(conn, "error", "failed to initialize terminal")
		_ = conn.Close()
		return
	}
	serverTerminal.Bridge(conn, terminalHandler)
	_ = conn.WriteJSON(response.Success("success"))
}

func parseServerID(value string) (uint, error) {
	return utils.StringToUint(value)
}
