package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/monitor/api/res"
	"squirrel-dev/internal/squ-apiserver/module/monitor/application"
	"squirrel-dev/internal/squ-apiserver/module/monitor/domain"
	"squirrel-dev/pkg/utils"
)

var errInvalidConfig = errors.New("invalid monitor configuration")

func (h *Handler) writeForServer(c *gin.Context, operation func(uint) (domain.Result, error)) {
	rawID := c.Param("serverId")
	id, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse monitor server ID", zap.String("raw_server_id", rawID), zap.Error(err))
		writeError(c, errInvalidConfig)
		return
	}
	result, err := operation(id)
	writeResult(c, result, err)
}

func writeResult(c *gin.Context, data domain.Result, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(data))
}

func writeError(c *gin.Context, err error) {
	code := res.ErrMonitorFailed
	switch {
	case errors.Is(err, errInvalidConfig):
		code = res.ErrInvalidMonitorConfig
	case errors.Is(err, application.ErrServerNotFound):
		code = res.ErrServerNotFound
	}
	c.JSON(http.StatusOK, response.Error(code))
}
