package monitor

import (
	"net/http"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/handler/monitor/res"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StatsHandler 查询系统统计数据
func StatsHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := service.GetStats()
		c.JSON(http.StatusOK, resp)
	}
}

// DiskIOHandler 查询单一磁盘IO
func DiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		device := c.Param("device")
		if device == "" {
			zap.L().Warn("Device parameter is empty")
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}
		resp := service.GetDiskIO(device)
		c.JSON(http.StatusOK, resp)
	}
}

// AllDiskIOHandler 查询全磁盘IO
func AllDiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := service.GetAllDiskIO()
		c.JSON(http.StatusOK, resp)
	}
}

// NetIOHandler 查询单一网卡流量
func NetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		interfaceName := c.Param("interface")
		if interfaceName == "" {
			zap.L().Warn("Interface parameter is empty")
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}
		resp := service.GetNetIO(interfaceName)
		c.JSON(http.StatusOK, resp)
	}
}

// AllNetIOHandler 查询全网卡流量
func AllNetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := service.GetAllNetIO()
		c.JSON(http.StatusOK, resp)
	}
}

// BaseMonitorRangeHandler 按时间范围查询基础监控数据
func BaseMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeRange := c.Query("range")
		if timeRange == "" {
			timeRange = "1h"
		}
		resp := service.GetBaseMonitorByRange(timeRange)
		c.JSON(http.StatusOK, resp)
	}
}

// DiskIOMonitorRangeHandler 按时间范围查询磁盘IO监控数据
func DiskIOMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeRange := c.Query("range")
		if timeRange == "" {
			timeRange = "1h"
		}
		resp := service.GetDiskIOMonitorByRange(timeRange)
		c.JSON(http.StatusOK, resp)
	}
}

// DiskUsageMonitorRangeHandler 按时间范围查询磁盘使用监控数据
func DiskUsageMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeRange := c.Query("range")
		if timeRange == "" {
			timeRange = "1h"
		}
		resp := service.GetDiskUsageMonitorByRange(timeRange)
		c.JSON(http.StatusOK, resp)
	}
}

// NetworkMonitorRangeHandler 按时间范围查询网络监控数据
func NetworkMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeRange := c.Query("range")
		if timeRange == "" {
			timeRange = "1h"
		}
		resp := service.GetNetworkMonitorByRange(timeRange)
		c.JSON(http.StatusOK, resp)
	}
}
