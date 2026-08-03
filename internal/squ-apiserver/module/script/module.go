package script

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/module/script/api"
	"squirrel-dev/internal/squ-apiserver/module/script/api/res"
	"squirrel-dev/internal/squ-apiserver/module/script/application"
	"squirrel-dev/internal/squ-apiserver/module/script/infra"
	serverInfra "squirrel-dev/internal/squ-apiserver/module/server/infra"
)

func BuildHandler(conf *config.Config, db *gorm.DB) *api.Handler {
	service := application.NewService(
		infra.NewRepository(db),
		infra.NewServerReader(serverInfra.NewRepository(db)),
		infra.NewAgentClient(conf),
		infra.IDGenerator{},
	)
	return api.NewHandler(service)
}

func RegisterHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	api.RegisterRoutes(group, BuildHandler(conf, db))
}

func RegisterAgentHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	api.RegisterAgentRoutes(group, BuildHandler(conf, db))
}

func MigrateScripts(db *gorm.DB) error  { return infra.MigrateScripts(db) }
func RollbackScripts(db *gorm.DB) error { return infra.RollbackScripts(db) }
func MigrateResults(db *gorm.DB) error  { return infra.MigrateResults(db) }
func RollbackResults(db *gorm.DB) error { return infra.RollbackResults(db) }
