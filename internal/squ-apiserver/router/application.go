package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/application"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"

	applicationRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

func Application(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	res.RegisterCode()

	service := application.New(
		conf,
		applicationRepository.New(db.GetDB()),
		serverRepository.New(db.GetDB()),
	)
	group.GET("/application", application.ListHandler(service))
	group.GET("/application/:id", application.GetHandler(service))
	group.DELETE("/application/:id", application.DeleteHandler(service))
	group.POST("/application", application.AddHandler(service))
	group.POST("/application/:id", application.UpdateHandler(service))
}
