package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/server", handler.List)
	group.GET("/server/:id", handler.Get)
	group.DELETE("/server/:id", handler.Delete)
	group.POST("/server", handler.Add)
	group.POST("/server/:id", handler.Update)
	group.POST("/server/check", handler.CheckAgent)
	group.POST("/ssh/test/:id", handler.TestSSH)
}

// RegisterTerminalRoute registers the WebSocket endpoint separately because it
// authenticates with the first WebSocket message rather than an HTTP header.
func RegisterTerminalRoute(group *gin.RouterGroup, handler *Handler) {
	group.GET("/ws/server/:id", handler.Terminal)
}
