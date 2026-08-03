package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/application", handler.List)
	group.GET("/application/:id", handler.Get)
	group.DELETE("/application/:id", handler.Delete)
	group.POST("/application", handler.Add)
	group.POST("/application/:id", handler.Update)
}
