package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/monitor"
	"squirrel-dev/internal/squ-agent/handler/monitor/res"
	monitorRepository "squirrel-dev/internal/squ-agent/repository/monitor"
	"squirrel-dev/pkg/collector"
)

func MonitorRouter(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	res.RegisterCode()

	// 创建收集器工厂并注册收集器
	factory := collector.NewCollectorFactory()
	factory.Register(collector.NewCPUCollector())
	factory.Register(collector.NewMemoryCollector())
	factory.Register(collector.NewDiskCollector())
	factory.Register(collector.NewIOCollector())
	factory.Register(collector.NewProcessCollector())

	// 创建monitor服务实例
	service := monitor.Monitor{
		Config:     conf,
		Repository: monitorRepository.New(db.GetDB()),
		Factory:    factory,
	}

	// 注册路由
	group.GET("/monitor/stats", monitor.StatsHandler(&service))
	group.GET("/monitor/stats/io/:device", monitor.DiskIOHandler(&service))
	group.GET("/monitor/stats/io/all", monitor.AllDiskIOHandler(&service))
	group.GET("/monitor/stats/net/:interface", monitor.NetIOHandler(&service))
	group.GET("/monitor/stats/net/all", monitor.AllNetIOHandler(&service))
	group.GET("/monitor/base/:page/:count", monitor.BaseMonitorPageHandler(&service))
	group.GET("/monitor/disk/:page/:count", monitor.DiskIOMonitorPageHandler(&service))
	group.GET("/monitor/net/:page/:count", monitor.NetworkMonitorPageHandler(&service))
}
