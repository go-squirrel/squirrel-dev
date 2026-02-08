package application

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := service.List()
		c.JSON(http.StatusOK, resp)
	}
}

func GetHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
			return
		}
		resp := service.Get(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
			return
		}
		resp := service.Delete(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func AddHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Application{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
			return
		}
		resp := service.Add(request)
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
			return
		}
		request := req.Application{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidApplicationConfig))
			return
		}
		request.ID = idUint
		resp := service.Update(request)
		c.JSON(http.StatusOK, resp)
	}
}
