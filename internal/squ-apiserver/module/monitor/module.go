package monitor

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/module/monitor/api"
	"squirrel-dev/internal/squ-apiserver/module/monitor/api/res"
	"squirrel-dev/internal/squ-apiserver/module/monitor/application"
	"squirrel-dev/internal/squ-apiserver/module/monitor/infra"
	serverInfra "squirrel-dev/internal/squ-apiserver/module/server/infra"
)

func RegisterHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	service := application.NewService(
		infra.NewServerReader(serverInfra.NewRepository(db)),
		infra.NewAgentClient(conf),
	)
	api.RegisterRoutes(group, api.NewHandler(service))
}
