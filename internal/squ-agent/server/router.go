package server

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/handler/health"
	"squirrel-dev/internal/squ-agent/router"
)

// SetupRouter 初始化gin入口，路由信息
func (s *Server) SetupRouter() {
	// 客户端通过变量进行传递
	v1Router := s.Gin.Group("/api/v1")
	healthRouter(v1Router)
	router.MonitorRouter(v1Router, s.Config, s.MonitorDB)
	router.Application(v1Router, s.Config, s.AppDB)
	router.Script(v1Router, s.Config, s.ScriptTaskDB)
	router.Server(v1Router, s.Config)
	router.Config(v1Router, s.Config, s.AgentDB)
}

func healthRouter(group *gin.RouterGroup) {
	group.GET("/health", health.HealthHandler())
}
