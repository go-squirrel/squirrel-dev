package monitor

import (
	"net/http"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/monitor/res"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StatsHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		resp := service.GetStats(serverIdUint)
		c.JSON(http.StatusOK, resp)
	}
}

func DiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		device := c.Param("device")
		if device == "" {
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		resp := service.GetDiskIO(serverIdUint, device)
		c.JSON(http.StatusOK, resp)
	}
}

func AllDiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		resp := service.GetAllDiskIO(serverIdUint)
		c.JSON(http.StatusOK, resp)
	}
}

func NetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		interfaceName := c.Param("interface")
		if interfaceName == "" {
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		resp := service.GetNetIO(serverIdUint, interfaceName)
		c.JSON(http.StatusOK, resp)
	}
}

func AllNetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		resp := service.GetAllNetIO(serverIdUint)
		c.JSON(http.StatusOK, resp)
	}
}

func BaseMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		timeRange := c.DefaultQuery("range", "1h")
		resp := service.GetBaseMonitorByRange(serverIdUint, timeRange)
		c.JSON(http.StatusOK, resp)
	}
}

func DiskIOMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		timeRange := c.DefaultQuery("range", "1h")
		resp := service.GetDiskIOMonitorByRange(serverIdUint, timeRange)
		c.JSON(http.StatusOK, resp)
	}
}

func DiskUsageMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		timeRange := c.DefaultQuery("range", "1h")
		resp := service.GetDiskUsageMonitorByRange(serverIdUint, timeRange)
		c.JSON(http.StatusOK, resp)
	}
}

func NetworkMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.L().Warn("failed to convert serverId to uint",
				zap.String("server_id", serverId),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}
		timeRange := c.DefaultQuery("range", "1h")
		resp := service.GetNetworkMonitorByRange(serverIdUint, timeRange)
		c.JSON(http.StatusOK, resp)
	}
}
