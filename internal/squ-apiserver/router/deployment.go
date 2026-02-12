package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/deployment"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"

	applicationRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	deploymentRepository "squirrel-dev/internal/squ-apiserver/repository/deployment"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

func Deployment(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	res.RegisterCode()

	service := deployment.New(
		conf,
		deploymentRepository.New(db.GetDB()),
		applicationRepository.New(db.GetDB()),
		serverRepository.New(db.GetDB()),
	)
	group.GET("/deployment", deployment.ListHandler(service))
	group.POST("/deployment/deploy/:id", deployment.DeployHandler(service))
	group.GET("/deployment/:id/servers", deployment.ListServersHandler(service))
	group.DELETE("/deployment/deploy/:id", deployment.UndeployHandler(service))
	group.POST("/deployment/stop/:id", deployment.StopHandler(service))
	group.POST("/deployment/start/:id", deployment.StartHandler(service))
	group.POST("/deployment/redeploy/:id", deployment.ReDeployHandler(service))
	group.POST("/deployment/report", deployment.ReportStatusHandler(service))
}
