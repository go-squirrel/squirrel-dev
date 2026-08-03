package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/deployment", handler.List)
	group.POST("/deployment/:id", handler.Update)
	group.POST("/deployment/deploy/:id", handler.Deploy)
	group.GET("/deployment/:id/servers", handler.ListServers)
	group.DELETE("/deployment/deploy/:id", handler.Undeploy)
	group.POST("/deployment/stop/:id", handler.Stop)
	group.POST("/deployment/start/:id", handler.Start)
	group.POST("/deployment/redeploy/:id", handler.ReDeploy)
}

func RegisterAgentRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/deployment/report", handler.ReportStatus)
}
