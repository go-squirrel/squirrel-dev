package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/config", handler.List)
	group.GET("/config/:id", handler.Get)
	group.DELETE("/config/:id", handler.Delete)
	group.POST("/config", handler.Save)
}
