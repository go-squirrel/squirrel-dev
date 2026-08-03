package application

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/application/api"
	"squirrel-dev/internal/squ-apiserver/module/application/api/res"
	appService "squirrel-dev/internal/squ-apiserver/module/application/application"
	"squirrel-dev/internal/squ-apiserver/module/application/infra"
)

func RegisterHTTP(group *gin.RouterGroup, db *gorm.DB) {
	res.RegisterCode()
	api.RegisterRoutes(group, api.NewHandler(appService.NewService(infra.NewRepository(db))))
}

func Migrate(db *gorm.DB) error { return infra.Migrate(db) }

func Rollback(db *gorm.DB) error { return infra.Rollback(db) }
