package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/script/api/res"
	"squirrel-dev/internal/squ-agent/module/script/application"
)

func bindRequest[T any](c *gin.Context) (T, bool) {
	var value T
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
		return value, false
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
	code := response.ErrSQL
	switch {
	case errors.Is(err, application.ErrAlreadyRunning):
		code = res.ErrScriptAlreadyRunning
	case errors.Is(err, gorm.ErrRecordNotFound):
		code = response.ErrSQLNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		code = response.ErrDuplicatedKey
	}
	c.JSON(http.StatusOK, response.Error(code))
}
