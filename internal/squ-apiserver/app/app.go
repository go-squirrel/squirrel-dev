package app

import (
	"time"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/middleware/cors"
	"squirrel-dev/internal/pkg/middleware/log"
	"squirrel-dev/internal/squ-apiserver/config"
	staticServer "squirrel-dev/internal/squ-apiserver/server"
)

// App 是应用级装配对象。
// 它负责承接配置、基础设施对象以及启动流程。
type App struct {
	Config *config.Config
	Gin    *gin.Engine
	Log    *log.Client
	DB     database.DB
}

func New() *App {
	return &App{}
}

// Run 启动整个应用。
func (a *App) Run() error {
	c := cors.New(cors.Config{
		AllowOrigins:     a.Config.Server.Origins,
		AllowMethods:     a.Config.Server.Methods,
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	})

	// Keep the legacy apiserver middleware behavior during the structural
	// migration. Logging and recovery were intentionally disabled there.
	staticServer.RegisterStatic(a.Gin)
	a.Gin.Use(c)

	a.runMigrations()
	a.registerHTTPRoutes()

	if err := a.Gin.Run(a.Config.Server.Bind + ":" + a.Config.Server.Port); err != nil {
		return err
	}
	return nil
}
