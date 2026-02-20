package server

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/middleware/mtls"
	"squirrel-dev/internal/squ-apiserver/handler/health"
	"squirrel-dev/internal/squ-apiserver/router"
)

// SetupRouter 初始化gin入口，路由信息
func (s *Server) SetupRouter() {
	// 客户端通过变量进行传递

	v1NoAuthRouter := s.Gin.Group("/api/v1")
	router.Auth(v1NoAuthRouter, s.Config, s.DB)
	// WebSocket 终端连接不需要 JWT 中间件认证，使用消息认证
	router.ServerTerminalNoAuth(v1NoAuthRouter, s.Config, s.DB)

	v1Router := s.Gin.Group("/api/v1")
	//v1Router.Use(jwt.JWTAuth(s.Config.Auth.Jwt.SigningKey))
	healthRouter(v1Router)
	router.Server(v1Router, s.Config, s.DB)
	router.Application(v1Router, s.Config, s.DB)
	router.Deployment(v1Router, s.Config, s.DB)
	router.Config(v1Router, s.Config, s.DB)
	router.AppStore(v1Router, s.Config, s.DB)
	router.Script(v1Router, s.Config, s.DB)
	router.Monitor(v1Router, s.Config, s.DB)
	s.SetupAgentRouter()
}

func (s *Server) SetupAgentRouter() {
	v1Router := s.Gin.Group("/api/v1")

	// Agent 路由使用 mTLS 双向认证
	if s.Config.MTLS.Enabled {
		v1Router.Use(mtls.MTLSAuthWithVerify(s.Config.MTLS.AllowedCNs))
	}

	router.AgentScript(v1Router, s.Config, s.DB)
	router.AgentDeployment(v1Router, s.Config, s.DB)
}

func healthRouter(group *gin.RouterGroup) {
	group.GET("/health", health.HealthHandler())
}
