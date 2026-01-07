package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/handler/monitor"
)

func MonitorRouter(group *gin.RouterGroup) {
	group.GET("/monitor", monitor.ListHandler())
}
