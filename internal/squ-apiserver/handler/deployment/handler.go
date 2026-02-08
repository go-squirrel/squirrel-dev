package deployment

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/req"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"
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
				c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
				return
			}
		}

		resp := service.List(serverIDUint)
		c.JSON(http.StatusOK, resp)
	}
}

func DeployHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}
		request := req.DeployApplication{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}
		request.ApplicationID = idUint
		resp := service.Deploy(request)
		c.JSON(http.StatusOK, resp)
	}
}

func ListServersHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}
		resp := service.ListServers(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func UndeployHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}

		resp := service.Undeploy(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func StopHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}

		resp := service.Stop(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func StartHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}

		resp := service.Start(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func ReportStatusHandler(service *Deployment) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.ReportApplicationStatus{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
			return
		}
		resp := service.ReportStatus(request)
		c.JSON(http.StatusOK, resp)
	}
}
