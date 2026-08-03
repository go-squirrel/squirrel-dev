package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/monitor/stats", handler.Stats)
	group.GET("/monitor/stats/io/:device", handler.DiskIO)
	group.GET("/monitor/stats/io/all", handler.AllDiskIO)
	group.GET("/monitor/stats/net/:interface", handler.NetIO)
	group.GET("/monitor/stats/net/all", handler.AllNetIO)
	group.GET("/monitor/base", handler.BaseByRange)
	group.GET("/monitor/disk", handler.DiskIOByRange)
	group.GET("/monitor/disk-usage", handler.DiskUsageByRange)
	group.GET("/monitor/net", handler.NetworkByRange)
}
