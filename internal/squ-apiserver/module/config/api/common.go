package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/config/api/req"
	"squirrel-dev/internal/squ-apiserver/module/config/api/res"
	"squirrel-dev/pkg/utils"
)

func bindConfig(c *gin.Context) (req.Config, bool) {
	var value req.Config
	if err := c.ShouldBindJSON(&value); err != nil {
		zap.L().Warn("failed to bind config request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigValue))
		return req.Config{}, false
	}
	return value, true
}

func configID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	id, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse config ID", zap.String("raw_config_id", rawID), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigKey))
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

func writeError(c *gin.Context, err error) {
	code := res.ErrConfigUpdateFailed
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		code = res.ErrConfigNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		code = res.ErrConfigKeyAlreadyExists
	}
	c.JSON(http.StatusOK, response.Error(code))
}
