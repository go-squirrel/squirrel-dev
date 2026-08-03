package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/script/api/res"
	"squirrel-dev/internal/squ-apiserver/module/script/application"
	"squirrel-dev/pkg/utils"
)

func bindRequest[T any](c *gin.Context) (T, bool) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("failed to bind script request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
		return request, false
	}
	return request, true
}

func scriptID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	id, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse script ID", zap.String("raw_script_id", rawID), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
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
	code := res.ErrInvalidScriptContent
	switch {
	case errors.Is(err, application.ErrNotFound):
		code = res.ErrScriptNotFound
	case errors.Is(err, application.ErrDuplicate):
		code = res.ErrDuplicateScript
	case errors.Is(err, application.ErrExecuteFailed):
		code = res.ErrScriptExecutionFailed
	case errors.Is(err, application.ErrServerNotFound):
		code = res.ErrServerNotFound
	}
	c.JSON(http.StatusOK, response.Error(code))
}
