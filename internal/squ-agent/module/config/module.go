package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/module/config/api"
	"squirrel-dev/internal/squ-agent/module/config/application"
	"squirrel-dev/internal/squ-agent/module/config/infra"
)

func RegisterHTTP(group *gin.RouterGroup, db *gorm.DB) {
	service := application.NewService(infra.NewRepository(db))
	api.RegisterRoutes(group, api.NewHandler(service))
}

func RegisterMigrations(registry *migration.MigrationRegistry) {
	infra.RegisterMigrations(registry)
}
