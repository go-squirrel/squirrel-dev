package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/scripts", handler.List)
	group.GET("/scripts/:id", handler.Get)
	group.DELETE("/scripts/:id", handler.Delete)
	group.POST("/scripts", handler.Add)
	group.POST("/scripts/:id", handler.Update)
	group.POST("/scripts/execute", handler.Execute)
	group.GET("/scripts/:id/results", handler.ListResults)
}

func RegisterAgentRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/scripts/receive-result", handler.ReceiveResult)
}
