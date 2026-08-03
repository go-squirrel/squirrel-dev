package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/config/api/req"
	"squirrel-dev/pkg/utils"
)

func bindConfig(c *gin.Context) (req.Config, bool) {
	var value req.Config
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
		return req.Config{}, false
	}
	return value, true
}

func configID(c *gin.Context) (uint, bool) {
	value, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
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
	code := response.ErrSQL
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		code = response.ErrSQLNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		code = response.ErrDuplicatedKey
	}
	c.JSON(http.StatusOK, response.Error(code))
}
