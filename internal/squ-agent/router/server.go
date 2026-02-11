package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/server"
	"squirrel-dev/internal/squ-agent/handler/server/res"
	"squirrel-dev/pkg/collector"
)

func Server(group *gin.RouterGroup, conf *config.Config) {
	res.RegisterCode()

	// 创建收集器工厂并注册收集器
	factory := collector.NewCollectorFactory()
	factory.Register(collector.NewHostCollector())

	// 创建server服务实例
	service := server.New(conf, factory)

	// 注册路由
	group.GET("/server/info", server.InfoHandler(service))
}
