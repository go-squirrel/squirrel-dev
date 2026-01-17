package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/application"

	applicationRepository "squirrel-dev/internal/squ-agent/repository/application"
)

func Application(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	service := application.Application{
		Config:     conf,
		Repository: applicationRepository.New(db.GetDB()),
	}
	group.GET("/application", application.ListHandler(&service))
	group.GET("/application/:id", application.GetHandler(&service))
	group.DELETE("/application/:id", application.DeleteHandler(&service))
	group.POST("/application", application.AddHandler(&service))
	group.POST("/application/:id", application.UpdateHandler(&service))
}
