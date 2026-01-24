package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/monitor"
	"squirrel-dev/internal/squ-apiserver/handler/monitor/res"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

func Monitor(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	res.RegisterCode()

	service := monitor.New(
		conf,
		serverRepository.New(db.GetDB()),
	)
	group.GET("/monitor/stats/:serverId", monitor.StatsHandler(service))
	group.GET("/monitor/stats/io/:serverId/:device", monitor.DiskIOHandler(service))
	group.GET("/monitor/stats/io/:serverId/all", monitor.AllDiskIOHandler(service))
	group.GET("/monitor/stats/net/:serverId/:interface", monitor.NetIOHandler(service))
	group.GET("/monitor/stats/net/:serverId/all", monitor.AllNetIOHandler(service))
	group.GET("/monitor/base/:serverId/:page/:count", monitor.BaseMonitorPageHandler(service))
	group.GET("/monitor/disk/:serverId/:page/:count", monitor.DiskIOMonitorPageHandler(service))
	group.GET("/monitor/net/:serverId/:page/:count", monitor.NetworkMonitorPageHandler(service))
}
