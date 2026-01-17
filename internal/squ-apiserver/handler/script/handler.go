package script

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := service.List()
		c.JSON(http.StatusOK, res)
	}
}

func GetHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Get(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func DeleteHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Delete(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func AddHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Script{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Add(request)
		c.JSON(http.StatusOK, res)
	}
}

func UpdateHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request := req.Script{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request.ID = idUint
		res := service.Update(request)
		c.JSON(http.StatusOK, res)
	}
}
