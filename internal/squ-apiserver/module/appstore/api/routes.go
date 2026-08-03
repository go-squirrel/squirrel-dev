package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/app-store", handler.List)
	group.GET("/app-store/:id", handler.Get)
	group.DELETE("/app-store/:id", handler.Delete)
	group.POST("/app-store", handler.Add)
	group.POST("/app-store/:id", handler.Update)
}
