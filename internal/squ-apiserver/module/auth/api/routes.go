package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {

}

func NoAuthRegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/login", handler.Login)
}
