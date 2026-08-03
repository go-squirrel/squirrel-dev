package health

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/health/api"
)

func RegisterHTTP(rg *gin.RouterGroup) {
	api.RegisterRoutes(rg, api.NewHandler())
}
