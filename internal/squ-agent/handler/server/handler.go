package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InfoHandler 获取服务器信息
func InfoHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := service.GetInfo()
		c.JSON(http.StatusOK, resp)
	}
}
