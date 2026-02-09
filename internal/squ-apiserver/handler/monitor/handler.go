package monitor

import (
	"net/http"
	"strconv"

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

func BaseMonitorPageHandler(service *Monitor) func(c *gin.Context) {
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
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("invalid page parameter",
				zap.String("page", pageStr),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("invalid count parameter",
				zap.String("count", countStr),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		resp := service.GetBaseMonitorPage(serverIdUint, page, count)
		c.JSON(http.StatusOK, resp)
	}
}

func DiskIOMonitorPageHandler(service *Monitor) func(c *gin.Context) {
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
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("invalid page parameter",
				zap.String("page", pageStr),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("invalid count parameter",
				zap.String("count", countStr),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		resp := service.GetDiskIOMonitorPage(serverIdUint, page, count)
		c.JSON(http.StatusOK, resp)
	}
}

func NetworkMonitorPageHandler(service *Monitor) func(c *gin.Context) {
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
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			zap.L().Warn("invalid page parameter",
				zap.String("page", pageStr),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			zap.L().Warn("invalid count parameter",
				zap.String("count", countStr),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
			return
		}

		resp := service.GetNetworkMonitorPage(serverIdUint, page, count)
		c.JSON(http.StatusOK, resp)
	}
}
