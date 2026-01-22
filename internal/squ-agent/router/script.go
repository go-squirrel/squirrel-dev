package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/script"
	scriptRes "squirrel-dev/internal/squ-agent/handler/script/res"

	scriptTaskRepo "squirrel-dev/internal/squ-agent/repository/script_task"
)

func Script(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	scriptRes.RegisterCode()

	service := script.New(
		conf,
		scriptTaskRepo.New(db.GetDB()),
	)

	// 执行脚本
	group.POST("/script/execute", script.ExecuteHandler(service))
}
