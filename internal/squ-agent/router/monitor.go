package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/monitor"
	"squirrel-dev/internal/squ-agent/handler/monitor/res"
	"squirrel-dev/pkg/collector"
)

func MonitorRouter(group *gin.RouterGroup, conf *config.Config) {
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
		Config:  conf,
		Factory: factory,
	}

	// 注册路由
	group.GET("/monitor/stats", monitor.StatsHandler(&service))
	group.GET("/monitor/stats/io/:device", monitor.DiskIOHandler(&service))
	group.GET("/monitor/stats/io/all", monitor.AllDiskIOHandler(&service))
	group.GET("/monitor/stats/net/:interface", monitor.NetIOHandler(&service))
	group.GET("/monitor/stats/net/all", monitor.AllNetIOHandler(&service))
}
