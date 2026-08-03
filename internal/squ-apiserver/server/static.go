package server

import (
	"embed"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	spaMiddleware "squirrel-dev/internal/pkg/middleware/static"
)

//go:embed all:dist
var staticData embed.FS

func RegisterStatic(engine *gin.Engine) {
	staticFiles, err := static.EmbedFolder(staticData, "dist")
	if err != nil {
		zap.L().Error("failed to embed static folder",
			zap.String("folder", "dist"),
			zap.Error(err),
		)
	}
	engine.Use(static.Serve("/", staticFiles))
	engine.Use(spaMiddleware.Default(staticData, "dist"))
}
