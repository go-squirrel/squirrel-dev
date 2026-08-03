package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/monitor/api/res"
	"squirrel-dev/internal/squ-agent/module/monitor/application"
)

var errInvalidParameter = errors.New("invalid monitor parameter")

func monitorDevice(c *gin.Context) (string, bool) {
	value := c.Param("device")
	if value == "" {
		zap.L().Warn("Device parameter is empty")
		writeError(c, errInvalidParameter)
		return "", false
	}
	return value, true
}

func networkInterface(c *gin.Context) (string, bool) {
	value := c.Param("interface")
	if value == "" {
		zap.L().Warn("Interface parameter is empty")
		writeError(c, errInvalidParameter)
		return "", false
	}
	return value, true
}

func monitorTimeRange(c *gin.Context) string {
	return c.DefaultQuery("range", "1h")
}

func writeResult(c *gin.Context, data any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(data))
}

func writeError(c *gin.Context, err error) {
	code := response.ErrSQL
	if errors.Is(err, errInvalidParameter) || errors.Is(err, application.ErrInvalidTimeRange) {
		zap.L().Warn("Invalid monitor parameter", zap.Error(err))
		code = res.ErrCodeParameter
	} else {
		zap.L().Error("Failed to get monitor data", zap.Error(err))
	}
	c.JSON(http.StatusOK, response.Error(code))
}
