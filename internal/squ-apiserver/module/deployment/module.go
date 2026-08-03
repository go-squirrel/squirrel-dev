package deployment

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/config"
	applicationInfra "squirrel-dev/internal/squ-apiserver/module/application/infra"
	"squirrel-dev/internal/squ-apiserver/module/deployment/api"
	"squirrel-dev/internal/squ-apiserver/module/deployment/api/res"
	"squirrel-dev/internal/squ-apiserver/module/deployment/application"
	"squirrel-dev/internal/squ-apiserver/module/deployment/infra"
	serverInfra "squirrel-dev/internal/squ-apiserver/module/server/infra"
)

func BuildHandler(conf *config.Config, db *gorm.DB) *api.Handler {
	service := application.NewService(
		infra.NewRepository(db),
		infra.NewApplicationReader(applicationInfra.NewRepository(db)),
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
func Migrate(db *gorm.DB) error  { return infra.Migrate(db) }
func Rollback(db *gorm.DB) error { return infra.Rollback(db) }
