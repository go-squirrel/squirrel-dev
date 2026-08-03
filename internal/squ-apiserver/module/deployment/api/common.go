package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/deployment/api/res"
	"squirrel-dev/internal/squ-apiserver/module/deployment/application"
	"squirrel-dev/pkg/utils"
)

func bindRequest[T any](c *gin.Context) (T, bool) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("failed to bind deployment request", zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
		return request, false
	}
	return request, true
}

func deploymentID(c *gin.Context) (uint, bool) {
	return pathID(c, "raw_deployment_id")
}

func applicationID(c *gin.Context) (uint, bool) {
	return pathID(c, "raw_application_id")
}

func pathID(c *gin.Context, field string) (uint, bool) {
	rawID := c.Param("id")
	id, err := utils.StringToUint(rawID)
	if err != nil {
		zap.L().Warn("failed to parse deployment path ID", zap.String(field, rawID), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
		return 0, false
	}
	return id, true
}

func deploymentServerID(c *gin.Context) (uint, bool) {
	value := c.Query("server_id")
	if value == "" {
		return 0, true
	}
	id, err := utils.StringToUint(value)
	if err != nil {
		zap.L().Warn("failed to parse deployment server ID", zap.String("raw_server_id", value), zap.Error(err))
		c.JSON(http.StatusOK, response.Error(res.ErrInvalidDeploymentConfig))
		return 0, false
	}
	return id, true
}

func writeResult(c *gin.Context, data any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Success(data))
}

func writeError(c *gin.Context, err error) {
	code := res.ErrCreateDeploymentRecordFailed
	switch {
	case errors.Is(err, application.ErrNotFound):
		code = res.ErrDeploymentNotFound
	case errors.Is(err, application.ErrAlreadyDeployed):
		code = res.ErrAlreadyDeployed
	case errors.Is(err, application.ErrApplicationMissing):
		code = res.ErrApplicationNotDeployed
	case errors.Is(err, application.ErrIDGeneration):
		code = res.ErrDeployIDGenerateFailed
	case errors.Is(err, application.ErrInvalidConfig):
		code = res.ErrInvalidDeploymentConfig
	case errors.Is(err, application.ErrContainerConflict):
		code = res.ErrComposeContainerNameConflict
	case errors.Is(err, application.ErrPortConflict):
		code = res.ErrComposePortConflict
	case errors.Is(err, application.ErrVolumeConflict):
		code = res.ErrComposeVolumeConflict
	case errors.Is(err, application.ErrNetworkConflict):
		code = res.ErrComposeNetworkConflict
	case errors.Is(err, application.ErrAgentDeploy):
		code = res.ErrAgentDeployFailed
	case errors.Is(err, application.ErrAgentDelete):
		code = res.ErrAgentDeleteFailed
	case errors.Is(err, application.ErrAgentStop):
		code = res.ErrAgentStopFailed
	case errors.Is(err, application.ErrAgentStart):
		code = res.ErrAgentStartFailed
	}
	c.JSON(http.StatusOK, response.Error(code))
}
