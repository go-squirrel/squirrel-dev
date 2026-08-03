package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/application", handler.List)
	group.GET("/application/:id", handler.Get)
	group.DELETE("/application/:id", handler.Delete)
	group.POST("/application", handler.Add)
	group.POST("/application/:id", handler.Update)
	group.POST("/application/start/:deployId", handler.Start)
	group.POST("/application/stop/:deployId", handler.Stop)
	group.POST("/application/delete/:deployId", handler.DeleteByDeployID)
}
