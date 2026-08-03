package application

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/module/application/api"
	"squirrel-dev/internal/squ-agent/module/application/api/res"
	appService "squirrel-dev/internal/squ-agent/module/application/application"
	appInfra "squirrel-dev/internal/squ-agent/module/application/infra"
	configInfra "squirrel-dev/internal/squ-agent/module/config/infra"
)

type Dependencies struct {
	Config  *config.Config
	AppDB   *gorm.DB
	AgentDB *gorm.DB
}

func RegisterHTTP(group *gin.RouterGroup, dependencies Dependencies) {
	res.RegisterCode()
	service := appService.NewService(
		appInfra.NewRepository(dependencies.AppDB),
		appInfra.NewConfigStore(configInfra.NewRepository(dependencies.AgentDB)),
		appInfra.NewComposeRuntime(dependencies.Config.Common.ComposePath),
	)
	api.RegisterRoutes(group, api.NewHandler(service))
}

func RegisterMigrations(registry *migration.MigrationRegistry) {
	appInfra.RegisterMigrations(registry)
}
