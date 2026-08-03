package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/server/api/res"
	"squirrel-dev/pkg/utils"
)

func bindRequest[T any](c *gin.Context) (T, bool) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("failed to bind server request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
		return request, false
	}
	return request, true
}

func serverID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	id, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse server ID", zap.String("raw_server_id", rawID), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
		return 0, false
	}
	return id, true
}

func writeResult(c *gin.Context, data any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Success(data))
}

func writeSSHResult(c *gin.Context, data any, err error) {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, response.Error(res.ErrSSHTestFailed))
		return
	}
	writeResult(c, data, err)
}

func writeError(c *gin.Context, err error) {
	code := res.ErrServerUpdateFailed
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		code = res.ErrServerNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		code = res.ErrServerAlreadyExists
	}
	c.JSON(http.StatusOK, response.Error(code))
}
