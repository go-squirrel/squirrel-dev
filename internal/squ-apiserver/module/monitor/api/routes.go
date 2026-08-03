package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/monitor/stats/:serverId", handler.Stats)
	group.GET("/monitor/stats/io/:serverId/:device", handler.DiskIO)
	group.GET("/monitor/stats/io/:serverId/all", handler.AllDiskIO)
	group.GET("/monitor/stats/net/:serverId/:interface", handler.NetIO)
	group.GET("/monitor/stats/net/:serverId/all", handler.AllNetIO)
	group.GET("/monitor/base/:serverId", handler.BaseRange)
	group.GET("/monitor/disk/:serverId", handler.DiskRange)
	group.GET("/monitor/disk-usage/:serverId", handler.DiskUsageRange)
	group.GET("/monitor/net/:serverId", handler.NetworkRange)
}
