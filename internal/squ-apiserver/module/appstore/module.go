package appstore

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/appstore/api"
	"squirrel-dev/internal/squ-apiserver/module/appstore/api/res"
	"squirrel-dev/internal/squ-apiserver/module/appstore/application"
	"squirrel-dev/internal/squ-apiserver/module/appstore/infra"
)

func RegisterHTTP(group *gin.RouterGroup, db *gorm.DB) {
	res.RegisterCode()
	api.RegisterRoutes(group, api.NewHandler(application.NewService(infra.NewRepository(db))))
}
func Migrate(db *gorm.DB) error  { return infra.Migrate(db) }
func Rollback(db *gorm.DB) error { return infra.Rollback(db) }
