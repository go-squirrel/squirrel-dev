package config

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/handler/config/req"
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
			zap.L().Warn("Failed to parse config id", zap.String("id", id), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Get(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func DeleteHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("Failed to parse config id", zap.String("id", id), zap.Error(err))
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Delete(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func SaveHandler(service *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Config{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("Failed to bind config request", zap.Error(err))
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Save(request)
		c.JSON(http.StatusOK, res)
	}
}
