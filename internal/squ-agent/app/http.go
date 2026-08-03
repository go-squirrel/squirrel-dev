package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
	applicationModule "squirrel-dev/internal/squ-agent/module/application"
	configModule "squirrel-dev/internal/squ-agent/module/config"
	healthModule "squirrel-dev/internal/squ-agent/module/health"
	monitorModule "squirrel-dev/internal/squ-agent/module/monitor"
	scriptModule "squirrel-dev/internal/squ-agent/module/script"
	serverModule "squirrel-dev/internal/squ-agent/module/server"
)

// registerHTTPRoutes 统一挂载所有 HTTP 路由。
// 应用层只负责装配和注册，不处理具体业务逻辑。
func (a *App) registerHTTPRoutes() {
	a.Gin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Success("ok"))
	})

	v1 := a.Gin.Group("/api/v1")
	healthModule.RegisterHTTP(v1)
	serverModule.RegisterHTTP(v1, serverModule.Dependencies{})
	if a.AgentDB != nil {
		configModule.RegisterHTTP(v1, a.AgentDB.GetDB())
	}
	if a.MonitorDB != nil {
		monitorModule.RegisterHTTP(v1, monitorModule.Dependencies{
			Cache: a.Cache,
			DB:    a.MonitorDB.GetDB(),
		})
	}
	if a.Config != nil && a.AppDB != nil && a.AgentDB != nil {
		applicationModule.RegisterHTTP(v1, applicationModule.Dependencies{
			Config:  a.Config,
			AppDB:   a.AppDB.GetDB(),
			AgentDB: a.AgentDB.GetDB(),
		})
	}
	if a.ScriptTaskDB != nil {
		scriptModule.RegisterHTTP(v1, a.ScriptTaskDB.GetDB())
	}
}
