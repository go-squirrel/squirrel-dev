package router

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/server"
	serverRes "squirrel-dev/internal/squ-apiserver/handler/server/res"

	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"github.com/gin-gonic/gin"
)

func Server(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	serverRes.RegisterCode()
	service := server.Server{
		Config:     conf,
		Repository: serverRepository.New(db.GetDB()),
	}
	group.GET("/server", server.ListHandler(&service))
	group.GET("/server/:id", server.GetHandler(&service))
	group.DELETE("/server/:id", server.DeleteHandler(&service))
	group.POST("/server", server.AddHandler(&service))
	group.POST("/server/:id", server.UpdateHandler(&service))

	group.GET("/ws/server/:id", server.TerminalHandler(&service))
}
