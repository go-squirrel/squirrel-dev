package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/application/api/req"
	"squirrel-dev/internal/squ-agent/module/application/api/res"
	"squirrel-dev/internal/squ-agent/module/application/application"
	"squirrel-dev/pkg/utils"
)

func bindApplication(c *gin.Context) (req.Application, bool) {
	var value req.Application
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
		return req.Application{}, false
	}
	return value, true
}

func applicationID(c *gin.Context) (uint, bool) {
	value, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
		return 0, false
	}
	return value, true
}

func deploymentID(c *gin.Context) (uint64, bool) {
	value, err := utils.StringToUint64(c.Param("deployId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
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

func writeApplicationResult(c *gin.Context, result application.Result) {
	writeResult(c, result.Data, result.Err)
}

func writeError(c *gin.Context, err error) {
	code := response.ErrSQL
	switch {
	case errors.Is(err, application.ErrDockerNotInstalled):
		code = res.ErrDockerNotInstalled
	case errors.Is(err, application.ErrComposeNotFound):
		code = res.ErrComposeNotFound
	case errors.Is(err, application.ErrComposeStart):
		code = res.ErrComposeStart
	case errors.Is(err, application.ErrComposeCreate):
		code = res.ErrComposeCreate
	case errors.Is(err, application.ErrComposeStop):
		code = res.ErrComposeStop
	case err == gorm.ErrRecordNotFound:
		code = response.ErrSQLNotFound
	case err == gorm.ErrDuplicatedKey:
		code = response.ErrDuplicatedKey
	}
	c.JSON(http.StatusOK, response.Error(code))
}
