package config

import (
	"fmt"
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/auth/req"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoginHandler(service *Config) func(c *gin.Context) {

	return func(c *gin.Context) {

		request := req.Request{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			fmt.Println(err)
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}

		res := service.List(request)

		c.JSON(http.StatusOK, res)
	}
}
