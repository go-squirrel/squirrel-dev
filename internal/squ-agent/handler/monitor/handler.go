package monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		service := Monitor{}
		res := service.Status()
		c.JSON(http.StatusOK, res)
	}
}
