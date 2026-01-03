package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		service := Health{}
		res := service.Status()
		c.JSON(http.StatusOK, res)
	}
}

