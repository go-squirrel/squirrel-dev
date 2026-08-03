package script

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/module/script/api"
	"squirrel-dev/internal/squ-agent/module/script/api/res"
	"squirrel-dev/internal/squ-agent/module/script/application"
	"squirrel-dev/internal/squ-agent/module/script/infra"
)

func RegisterHTTP(group *gin.RouterGroup, db *gorm.DB) {
	res.RegisterCode()
	service := application.NewService(infra.NewRepository(db), infra.NewShellExecutor())
	api.RegisterRoutes(group, api.NewHandler(service))
}

func RegisterMigrations(registry *migration.MigrationRegistry) {
	infra.RegisterMigrations(registry)
}
