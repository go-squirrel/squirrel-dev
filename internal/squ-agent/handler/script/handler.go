package script

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	scriptReq "squirrel-dev/internal/squ-agent/handler/script/req"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ExecuteHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := scriptReq.Script{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("Failed to bind script request", zap.Error(err))
			c.JSON(http.StatusOK, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Execute(request)
		c.JSON(http.StatusOK, res)
	}
}
