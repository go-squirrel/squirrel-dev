package monitor

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/cache"
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/module/monitor/api"
	"squirrel-dev/internal/squ-agent/module/monitor/api/res"
	"squirrel-dev/internal/squ-agent/module/monitor/application"
	"squirrel-dev/internal/squ-agent/module/monitor/domain"
	"squirrel-dev/internal/squ-agent/module/monitor/infra"
)

type Dependencies struct {
	Cache     cache.Cache
	DB        *gorm.DB
	Collector domain.Collector
}

func RegisterHTTP(group *gin.RouterGroup, dependencies Dependencies) {
	res.RegisterCode()
	metricsCollector := dependencies.Collector
	if metricsCollector == nil {
		metricsCollector = infra.NewDefaultCollector()
	}
	service := application.NewService(
		dependencies.Cache,
		infra.NewRepository(dependencies.DB),
		metricsCollector,
	)
	api.RegisterRoutes(group, api.NewHandler(service))
}

func RegisterMigrations(registry *migration.MigrationRegistry) {
	infra.RegisterMigrations(registry)
}
