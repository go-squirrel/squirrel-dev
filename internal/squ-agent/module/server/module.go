package server

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/server/api"
	serverRes "squirrel-dev/internal/squ-agent/module/server/api/res"
	"squirrel-dev/internal/squ-agent/module/server/application"
	"squirrel-dev/internal/squ-agent/module/server/infra"
	"squirrel-dev/pkg/collector"
)

type Dependencies struct {
	HostCollector collector.HostCollector
}

func RegisterHTTP(rg *gin.RouterGroup, dependencies Dependencies) {
	serverRes.RegisterCode()

	hostCollector := dependencies.HostCollector
	if hostCollector == nil {
		hostCollector = collector.NewHostCollector()
	}

	service := application.NewService(infra.NewHostInfoCollector(hostCollector))
	api.RegisterRoutes(rg, api.NewHandler(service))
}
