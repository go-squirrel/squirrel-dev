package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/application/api/req"
	"squirrel-dev/internal/squ-apiserver/module/application/api/res"
	appService "squirrel-dev/internal/squ-apiserver/module/application/application"
	"squirrel-dev/pkg/utils"
)

func bindApplication(c *gin.Context) (req.Application, bool) {
	var value req.Application
	if err := c.ShouldBindJSON(&value); err != nil {
		zap.L().Warn("failed to bind application request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
		return req.Application{}, false
	}
	return value, true
}

func applicationID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	value, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse application ID", zap.String("raw_application_id", rawID), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
		return 0, false
	}
	return value, true
}

func writeResult(c *gin.Context, data any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Success(data))
}

func writeError(c *gin.Context, err error) {
	code := res.ErrApplicationUpdateFailed
	switch {
	case errors.Is(err, appService.ErrInvalidYAML):
		code = res.ErrInvalidApplicationConfig
	case err == gorm.ErrRecordNotFound:
		code = res.ErrApplicationNotFound
	case err == gorm.ErrDuplicatedKey:
		code = res.ErrDuplicateApplication
	}
	c.JSON(http.StatusOK, response.Error(code))
}
