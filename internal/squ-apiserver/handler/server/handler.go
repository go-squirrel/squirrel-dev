package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListHandler(service *Server) func(c *gin.Context) {

	return func(c *gin.Context) {
		res := service.List()
		c.JSON(http.StatusOK, res)
	}
}
