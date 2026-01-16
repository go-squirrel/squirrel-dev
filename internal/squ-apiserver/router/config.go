package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	sysConfig "squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/config"

	configRepository "squirrel-dev/internal/squ-apiserver/repository/config"
)

func Config(group *gin.RouterGroup, conf *sysConfig.Config, db database.DB) {
	service := config.Config{
		Config:     conf,
		Repository: configRepository.New(db.GetDB()),
	}
	group.GET("/config", config.ListHandler(&service))
	group.GET("/config/:id", config.GetHandler(&service))
	group.DELETE("/config/:id", config.DeleteHandler(&service))
	group.POST("/config", config.AddHandler(&service))
	group.POST("/config/:id", config.UpdateHandler(&service))
}
