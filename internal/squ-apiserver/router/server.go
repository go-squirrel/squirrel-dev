package router

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server"
	serverRes "squirrel-dev/internal/squ-apiserver/handler/server/res"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"github.com/gin-gonic/gin"
)

func Server(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	serverRes.RegisterCode()

	service := server.New(conf, serverRepository.New(db.GetDB()))
	group.GET("/server", server.ListHandler(service))
	group.GET("/server/:id", server.GetHandler(service))
	group.DELETE("/server/:id", server.DeleteHandler(service))
	group.POST("/server", server.AddHandler(service))
	group.POST("/server/:id", server.UpdateHandler(service))
	group.POST("/server/check", server.CheckAgentHandler(service))
	group.POST("/ssh/test/:id", server.TestSSHHandler(service))
}

// ServerTerminalNoAuth 终端 WebSocket 路由（无 JWT 中间件，使用消息认证）
func ServerTerminalNoAuth(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	service := server.New(conf, serverRepository.New(db.GetDB()))
	group.GET("/ws/server/:id", server.TerminalHandler(service))
}
