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
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res, err := service.GetStats(serverIdUint)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func DiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		device := c.Param("device")
		if device == "" {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res, err := service.GetDiskIO(serverIdUint, device)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func AllDiskIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res, err := service.GetAllDiskIO(serverIdUint)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func NetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		interfaceName := c.Param("interface")
		if interfaceName == "" {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res, err := service.GetNetIO(serverIdUint, interfaceName)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func AllNetIOHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res, err := service.GetAllNetIO(serverIdUint)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func BaseMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res, err := service.GetBaseMonitorPage(serverIdUint, page, count)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func DiskIOMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res, err := service.GetDiskIOMonitorPage(serverIdUint, page, count)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func NetworkMonitorPageHandler(service *Monitor) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverId := c.Param("serverId")
		serverIdUint, err := utils.StringToUint(serverId)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		pageStr := c.Param("page")
		countStr := c.Param("count")

		if pageStr == "" || countStr == "" {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 || count > 100 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res, err := service.GetNetworkMonitorPage(serverIdUint, page, count)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrCodeParameter))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
