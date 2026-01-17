package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	scriptHandler "squirrel-dev/internal/squ-apiserver/handler/script"
	scriptRes "squirrel-dev/internal/squ-apiserver/handler/script/res"

	scriptRepository "squirrel-dev/internal/squ-apiserver/repository/script"
)

func Script(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	scriptRes.RegisterCode()

	service := scriptHandler.Script{
		Config:     conf,
		Repository: scriptRepository.New(db.GetDB()),
	}

	// 数据库中的脚本 CRUD 操作
	group.GET("/scripts", scriptHandler.ListHandler(&service))
	group.GET("/scripts/:id", scriptHandler.GetHandler(&service))
	group.DELETE("/scripts/:id", scriptHandler.DeleteHandler(&service))
	group.POST("/scripts", scriptHandler.AddHandler(&service))
	group.POST("/scripts/:id", scriptHandler.UpdateHandler(&service))
}
