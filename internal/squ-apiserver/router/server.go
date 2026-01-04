package router

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server"

	serverModel "squirrel-dev/internal/squ-apiserver/model/server"

	"github.com/gin-gonic/gin"
)

func Server(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	service := server.Server{
		Config:      conf,
		ModelClient: serverModel.New(db.GetDB()),
	}
	group.GET("/server", server.ListHandler(&service))
}
