package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/auth/api/req"
	"squirrel-dev/internal/squ-apiserver/module/auth/api/res"
	"squirrel-dev/internal/squ-apiserver/module/auth/application"
)

func bindLogin(c *gin.Context) (req.Request, bool) {
	var request req.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("failed to bind login request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidCredentials))
		return req.Request{}, false
	}
	return request, true
}

func writeResult(c *gin.Context, data any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Success(data))
}

func writeError(c *gin.Context, err error) {
	code := res.ErrAuthFailed
	switch {
	case errors.Is(err, application.ErrInvalidCredentials):
		code = res.ErrInvalidCredentials
	case errors.Is(err, application.ErrTokenGeneration):
		code = res.ErrTokenGenerateFailed
	}
	c.JSON(http.StatusOK, response.Error(code))
}
