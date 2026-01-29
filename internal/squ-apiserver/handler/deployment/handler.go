package deployment

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/req"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverID := c.Query("server_id")
		var serverIDUint uint = 0
		var err error

		if serverID != "" {
			serverIDUint, err = utils.StringToUint(serverID)
			if err != nil {
				zap.S().Warn(err)
				c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
				return
			}
		}

		res := service.List(serverIDUint)
		c.JSON(http.StatusOK, res)
	}
}

func DeployHandler(service *Deployment) func(c *gin.Context) {
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

func ListServersHandler(service *Deployment) func(c *gin.Context) {
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

func UndeployHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.Undeploy(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func StopHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.Stop(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func StartHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.Start(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func ReportStatusHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.ReportApplicationStatus{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.ReportStatus(request)
		c.JSON(http.StatusOK, res)
	}
}
