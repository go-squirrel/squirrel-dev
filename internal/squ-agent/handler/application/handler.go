package application

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/handler/application/req"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := service.List()
		c.JSON(http.StatusOK, res)
	}
}

func GetHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("Failed to parse application id", zap.String("id", id), zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Get(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func DeleteHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("Failed to parse application id", zap.String("id", id), zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Delete(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func AddHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Application{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("Failed to bind application request", zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Add(request)
		c.JSON(http.StatusOK, res)
	}
}

func UpdateHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("Failed to parse application id", zap.String("id", id), zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request := req.Application{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("Failed to bind application request", zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request.ID = idUint
		res := service.Update(request)
		c.JSON(http.StatusOK, res)
	}
}

func StopHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		deployID := c.Param("deployId")
		deployIDUint64, err := utils.StringToUint64(deployID)
		if err != nil {
			zap.L().Warn("Failed to parse deploy id", zap.String("deployId", deployID), zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Stop(deployIDUint64)
		c.JSON(http.StatusOK, res)
	}
}

func StartHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		deployID := c.Param("deployId")
		deployIDUint64, err := utils.StringToUint64(deployID)
		if err != nil {
			zap.L().Warn("Failed to parse deploy id", zap.String("deployId", deployID), zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Start(deployIDUint64)
		c.JSON(http.StatusOK, res)
	}
}

// DeleteByNameHandler 根据deployID删除应用（用于回滚）
func DeleteByNameHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		deployID := c.Param("deployId")
		deployIDUint64, err := utils.StringToUint64(deployID)
		if err != nil {
			zap.L().Warn("Failed to parse deploy id", zap.String("deployId", deployID), zap.Error(err))
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.DeleteByDeployID(deployIDUint64)
		c.JSON(http.StatusOK, res)
	}
}
