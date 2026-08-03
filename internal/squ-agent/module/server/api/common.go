package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
)

func writeResult(c *gin.Context, data any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Success(data))
}

func writeError(c *gin.Context, err error) {
	zap.L().Error("Failed to collect host info", zap.Error(err))
	// Preserve the legacy model.ReturnErrCode behavior for collector errors.
	c.JSON(http.StatusOK, response.Error(response.ErrSQL))
}
