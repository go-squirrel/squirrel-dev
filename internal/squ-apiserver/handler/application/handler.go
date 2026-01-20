package application

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
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
			zap.S().Warn(err)
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
			zap.S().Warn(err)
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
			zap.S().Warn(err)
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
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request := req.Application{}
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

func DeployHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request := req.DeployApplication{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request.ApplicationID = idUint
		res := service.Deploy(request)
		c.JSON(http.StatusOK, res)
	}
}

func ListServersHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.ListServers(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func UndeployHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		serverID := c.Param("serverId")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		serverIDUint, err := utils.StringToUint(serverID)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.Undeploy(idUint, serverIDUint)
		c.JSON(http.StatusOK, res)
	}
}

func StopHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		serverID := c.Param("serverId")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		serverIDUint, err := utils.StringToUint(serverID)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.Stop(idUint, serverIDUint)
		c.JSON(http.StatusOK, res)
	}
}

func StartHandler(service *Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		serverID := c.Param("serverId")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		serverIDUint, err := utils.StringToUint(serverID)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.Start(idUint, serverIDUint)
		c.JSON(http.StatusOK, res)
	}
}

