package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/module/auth/api"
	"squirrel-dev/internal/squ-apiserver/module/auth/api/res"
	"squirrel-dev/internal/squ-apiserver/module/auth/application"
	"squirrel-dev/internal/squ-apiserver/module/auth/infra"
)

func NoAuthRegisterHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	service := application.NewService(
		infra.NewVerifier(db),
		infra.NewTokenGenerator(conf.Auth.Jwt.SigningKey, conf.Auth.Jwt.Expired),
	)
	api.NoAuthRegisterRoutes(group, api.NewHandler(service))
}

func RegisterHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	service := application.NewService(
		infra.NewVerifier(db),
		infra.NewTokenGenerator(conf.Auth.Jwt.SigningKey, conf.Auth.Jwt.Expired),
	)
	api.RegisterRoutes(group, api.NewHandler(service))
}

func Migrate(db *gorm.DB) error  { return infra.Migrate(db) }
func Rollback(db *gorm.DB) error { return infra.Rollback(db) }
