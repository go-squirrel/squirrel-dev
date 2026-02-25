package monitor

import (
	"net/http"
	"strconv"

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

// BaseMonitorPageHandler 查询基础监控数据分页
func BaseMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			zap.L().Warn("Page or count parameter is empty")
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("Invalid page parameter", zap.String("page", pageStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("Invalid count parameter", zap.String("count", countStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		resp := service.GetBaseMonitorPage(page, count)
		c.JSON(http.StatusOK, resp)
	}
}

// DiskIOMonitorPageHandler 查询磁盘IO监控数据分页
func DiskIOMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			zap.L().Warn("Page or count parameter is empty")
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("Invalid page parameter", zap.String("page", pageStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("Invalid count parameter", zap.String("count", countStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		resp := service.GetDiskIOMonitorPage(page, count)
		c.JSON(http.StatusOK, resp)
	}
}

// NetworkMonitorPageHandler query network monitor data page
func NetworkMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			zap.L().Warn("Page or count parameter is empty")
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("Invalid page parameter", zap.String("page", pageStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("Invalid count parameter", zap.String("count", countStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		resp := service.GetNetworkMonitorPage(page, count)
		c.JSON(http.StatusOK, resp)
	}
}

// DiskUsageMonitorPageHandler query disk usage monitor data page
func DiskUsageMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			zap.L().Warn("Page or count parameter is empty")
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("Invalid page parameter", zap.String("page", pageStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("Invalid count parameter", zap.String("count", countStr), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(res.ErrCodeParameter))
			return
		}

		resp := service.GetDiskUsageMonitorPage(page, count)
		c.JSON(http.StatusOK, resp)
	}
}
