package server

import (
	"net/http"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/handler/server/res"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InfoHandler 获取服务器信息
func InfoHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := service.GetInfo()
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusInternalServerError, response.Error(res.ErrCodeSystem))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
