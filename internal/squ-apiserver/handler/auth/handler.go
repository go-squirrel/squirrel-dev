package auth

import (
	"fmt"
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/auth/req"
	"squirrel-dev/internal/squ-apiserver/handler/auth/res"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoginHandler(service *Auth) func(c *gin.Context) {

	return func(c *gin.Context) {

		request := req.Request{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			fmt.Println(err)
			zap.S().Warn(err)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidCredentials))
			return
		}

		resp := service.Login(request)

		c.JSON(http.StatusOK, resp)
	}
}
