package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/module/server/api"
	"squirrel-dev/internal/squ-apiserver/module/server/api/res"
	"squirrel-dev/internal/squ-apiserver/module/server/application"
	"squirrel-dev/internal/squ-apiserver/module/server/infra"
	"squirrel-dev/pkg/httpclient"
)

func buildHandler(conf *config.Config, db *gorm.DB) *api.Handler {
	service := application.NewService(
		infra.NewRepository(db),
		infra.NewAgentClient(conf, httpclient.NewClient(3*time.Second)),
		infra.NewSSHTester(),
	)
	return api.NewHandler(service, conf.Auth.Jwt.SigningKey)
}

func RegisterHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	api.RegisterRoutes(group, buildHandler(conf, db))
}

// RegisterTerminalHTTP keeps the WebSocket route outside the HTTP JWT
// middleware. The terminal handler validates the token sent by the client in
// the first WebSocket message.
func RegisterTerminalHTTP(group *gin.RouterGroup, conf *config.Config, db *gorm.DB) {
	res.RegisterCode()
	api.RegisterTerminalRoute(group, buildHandler(conf, db))
}

func Migrate(db *gorm.DB) error  { return infra.Migrate(db) }
func Rollback(db *gorm.DB) error { return infra.Rollback(db) }
