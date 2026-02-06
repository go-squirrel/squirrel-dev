package monitor

import (
	"net/http"
	"strconv"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StatsHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.GetStats(serverIdUint)
		c.JSON(http.StatusOK, res)
	}
}

func DiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		device := c.Param("device")
		if device == "" {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.GetDiskIO(serverIdUint, device)
		c.JSON(http.StatusOK, res)
	}
}

func AllDiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.GetAllDiskIO(serverIdUint)
		c.JSON(http.StatusOK, res)
	}
}

func NetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		interfaceName := c.Param("interface")
		if interfaceName == "" {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.GetNetIO(serverIdUint, interfaceName)
		c.JSON(http.StatusOK, res)
	}
}

func AllNetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.GetAllNetIO(serverIdUint)
		c.JSON(http.StatusOK, res)
	}
}

func BaseMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.GetBaseMonitorPage(serverIdUint, page, count)
		c.JSON(http.StatusOK, res)
	}
}

func DiskIOMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.GetDiskIOMonitorPage(serverIdUint, page, count)
		c.JSON(http.StatusOK, res)
	}
}

func NetworkMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.GetNetworkMonitorPage(serverIdUint, page, count)
		c.JSON(http.StatusOK, res)
	}
}
