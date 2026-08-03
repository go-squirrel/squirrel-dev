package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/appstore/api/req"
	"squirrel-dev/internal/squ-apiserver/module/appstore/api/res"
	"squirrel-dev/internal/squ-apiserver/module/appstore/application"
	"squirrel-dev/pkg/utils"
)

func bindAppStore(c *gin.Context) (req.AppStore, bool) {
	var value req.AppStore
	if err := c.ShouldBindJSON(&value); err != nil {
		zap.L().Warn("failed to bind app store request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
		return req.AppStore{}, false
	}
	return value, true
}

func appStoreID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	id, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse app store ID", zap.String("raw_app_store_id", rawID), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
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
	code := res.ErrAppStoreUpdateFailed
	switch {
	case errors.Is(err, application.ErrInvalidCompose):
		code = res.ErrInvalidComposeContent
	case errors.Is(err, application.ErrUnsupportedType):
		code = res.ErrUnsupportedAppType
	case err == gorm.ErrRecordNotFound:
		code = res.ErrAppStoreNotFound
	case err == gorm.ErrDuplicatedKey:
		code = res.ErrDuplicateAppStore
	}
	c.JSON(http.StatusOK, response.Error(code))
}
