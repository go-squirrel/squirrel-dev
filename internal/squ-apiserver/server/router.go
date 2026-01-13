package server

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-apiserver/handler/health"
	"squirrel-dev/internal/squ-apiserver/router"
)

// SetupRouter 初始化gin入口，路由信息
func (s *Server) SetupRouter() {
	// 客户端通过变量进行传递
	v1Router := s.Gin.Group("/api/v1")
	healthRouter(v1Router)
	router.Auth(v1Router, s.Config, s.DB)
	router.Server(v1Router, s.Config, s.DB)
	router.Application(v1Router, s.Config, s.DB)
	router.Config(v1Router, s.Config, s.DB)
}

func healthRouter(group *gin.RouterGroup) {
	group.GET("/health", health.HealthHandler())
}
