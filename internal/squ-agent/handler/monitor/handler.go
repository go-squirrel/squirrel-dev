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
		resp, err := service.GetStats()
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(res.ErrCodeSystem))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// DiskIOHandler 查询单一磁盘IO
func DiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		device := c.Param("device")
		if device == "" {
			c.JSON(http.StatusBadRequest, response.Error(res.ErrCodeParameter))
			return
		}
		resp, err := service.GetDiskIO(device)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(res.ErrCodeSystem))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// AllDiskIOHandler 查询全磁盘IO
func AllDiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := service.GetAllDiskIO()
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(res.ErrCodeSystem))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// NetIOHandler 查询单一网卡流量
func NetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		interfaceName := c.Param("interface")
		if interfaceName == "" {
			c.JSON(http.StatusBadRequest, response.Error(res.ErrCodeParameter))
			return
		}
		resp, err := service.GetNetIO(interfaceName)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(res.ErrCodeSystem))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// AllNetIOHandler 查询全网卡流量
func AllNetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := service.GetAllNetIO()
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(res.ErrCodeSystem))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
