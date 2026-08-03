package app

import (
	"net/http"

	"squirrel-dev/internal/pkg/jwt"
	"squirrel-dev/internal/pkg/middleware/mtls"
	"squirrel-dev/internal/pkg/response"
	applicationModule "squirrel-dev/internal/squ-apiserver/module/application"
	appstoreModule "squirrel-dev/internal/squ-apiserver/module/appstore"
	authModule "squirrel-dev/internal/squ-apiserver/module/auth"
	configModule "squirrel-dev/internal/squ-apiserver/module/config"
	deploymentModule "squirrel-dev/internal/squ-apiserver/module/deployment"
	monitorModule "squirrel-dev/internal/squ-apiserver/module/monitor"
	scriptModule "squirrel-dev/internal/squ-apiserver/module/script"
	serverModule "squirrel-dev/internal/squ-apiserver/module/server"

	"github.com/gin-gonic/gin"
)

// registerHTTPRoutes 统一挂载所有 HTTP 路由。
// 应用层只负责装配和注册，不处理具体业务逻辑。
func (a *App) registerHTTPRoutes() {
	a.Gin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Success("ok"))
	})

	v1 := a.Gin.Group("/api/v1")
	if a.Config != nil && a.DB != nil {
		authModule.NoAuthRegisterHTTP(v1, a.Config, a.DB.GetDB())
		// 与旧版一致：终端 WebSocket 不经过 HTTP JWT 中间件，而是在
		// WebSocket 建立后通过首条 auth 消息校验 token。
		serverModule.RegisterTerminalHTTP(v1, a.Config, a.DB.GetDB())

		v1Auth := a.Gin.Group("/api/v1")
		v1Auth.Use(jwt.JWTAuth(a.Config.Auth.Jwt.SigningKey))
		authModule.RegisterHTTP(v1Auth, a.Config, a.DB.GetDB())
		serverModule.RegisterHTTP(v1Auth, a.Config, a.DB.GetDB())
		configModule.RegisterHTTP(v1Auth, a.DB.GetDB())
		appstoreModule.RegisterHTTP(v1Auth, a.DB.GetDB())
		applicationModule.RegisterHTTP(v1Auth, a.DB.GetDB())
		deploymentModule.RegisterHTTP(v1Auth, a.Config, a.DB.GetDB())
		scriptModule.RegisterHTTP(v1Auth, a.Config, a.DB.GetDB())
		monitorModule.RegisterHTTP(v1Auth, a.Config, a.DB.GetDB())
	}
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Success("health"))
	})

	if a.Config != nil && a.DB != nil {
		agentV1 := a.Gin.Group("/api/v1")
		if a.Config.MTLS.Enabled {
			agentV1.Use(mtls.MTLSAuthWithVerify(a.Config.MTLS.AllowedCNs))
		}
		deploymentModule.RegisterAgentHTTP(agentV1, a.Config, a.DB.GetDB())
		scriptModule.RegisterAgentHTTP(agentV1, a.Config, a.DB.GetDB())
	}
}
