package config

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/config/req"
	"squirrel-dev/internal/squ-apiserver/handler/config/res"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := service.List()
		c.JSON(http.StatusOK, res)
	}
}

func GetHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigKey))
			return
		}
		resp := service.Get(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigKey))
			return
		}
		resp := service.Delete(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func AddHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Config{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigValue))
			return
		}
		resp := service.Add(request)
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigKey))
			return
		}
		request := req.Config{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidConfigValue))
			return
		}
		request.ID = idUint
		resp := service.Update(request)
		c.JSON(http.StatusOK, resp)
	}
}

