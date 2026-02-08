package app_store

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/req"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/res"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *AppStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := service.List()
		c.JSON(http.StatusOK, resp)
	}
}

func GetHandler(service *AppStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
			return
		}
		resp := service.Get(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteHandler(service *AppStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
			return
		}
		resp := service.Delete(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func AddHandler(service *AppStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.AppStore{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
			return
		}
		resp := service.Add(request)
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateHandler(service *AppStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
			return
		}
		request := req.AppStore{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidAppStoreConfig))
			return
		}
		request.ID = idUint
		resp := service.Update(request)
		c.JSON(http.StatusOK, resp)
	}
}
