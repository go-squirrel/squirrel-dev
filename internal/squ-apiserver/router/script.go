package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	scriptHandler "squirrel-dev/internal/squ-apiserver/handler/script"
	scriptRes "squirrel-dev/internal/squ-apiserver/handler/script/res"

	scriptRepository "squirrel-dev/internal/squ-apiserver/repository/script"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

func Script(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	scriptRes.RegisterCode()

	service := scriptHandler.New(
		conf,
		scriptRepository.New(db.GetDB()),
		serverRepository.New(db.GetDB()),
	)

	// 数据库中的脚本 CRUD 操作
	group.GET("/scripts", scriptHandler.ListHandler(service))
	group.GET("/scripts/:id", scriptHandler.GetHandler(service))
	group.DELETE("/scripts/:id", scriptHandler.DeleteHandler(service))
	group.POST("/scripts", scriptHandler.AddHandler(service))
	group.POST("/scripts/:id", scriptHandler.UpdateHandler(service))

	// 执行脚本相关操作
	group.POST("/scripts/execute", scriptHandler.ExecuteHandler(service))
	group.GET("/scripts/:id/results", scriptHandler.GetResultsHandler(service))
}

func AgentScript(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	scriptRes.RegisterCode()

	service := scriptHandler.New(
		conf,
		scriptRepository.New(db.GetDB()),
		serverRepository.New(db.GetDB()),
	)
	group.POST("/scripts/receive-result", scriptHandler.ReceiveResultHandler(service))
}
