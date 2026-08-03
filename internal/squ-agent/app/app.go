package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/cache"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/middleware/cors"
	"squirrel-dev/internal/pkg/middleware/log"
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/config"
	applicationModule "squirrel-dev/internal/squ-agent/module/application"
	configModule "squirrel-dev/internal/squ-agent/module/config"
	monitorModule "squirrel-dev/internal/squ-agent/module/monitor"
	scriptModule "squirrel-dev/internal/squ-agent/module/script"
)

// App 是应用级装配对象。
// 它负责承接配置、基础设施对象以及启动流程。
type App struct {
	Config       *config.Config
	Gin          *gin.Engine
	Log          *log.Client
	AgentDB      database.DB
	AppDB        database.DB
	MonitorDB    database.DB
	ScriptTaskDB database.DB
	Cache        cache.Cache
	Jobs         interface{ Start() error }
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

	a.Gin.Use(
		log.GinLogger(a.Log.Logger),
		log.GinRecovery(a.Log.Logger, true),
		c,
	)

	a.runMigrations()
	a.registerHTTPRoutes()
	if a.Jobs != nil {
		if err := a.Jobs.Start(); err != nil {
			zap.L().Warn("failed to start agent jobs", zap.Error(err))
		}
	}

	if err := a.Gin.Run(a.Config.Server.Bind + ":" + a.Config.Server.Port); err != nil {
		return err
	}
	return nil
}

// runMigrations preserves the legacy startup behavior: migration failures are
// logged, but they do not prevent the HTTP server from starting.
func (a *App) runMigrations() {
	a.runMigration("agent", a.AgentDB, a.buildAgentMigrationRegistry())
	a.runMigration("application", a.AppDB, a.buildApplicationMigrationRegistry())
	a.runMigration("script task", a.ScriptTaskDB, a.buildScriptTaskMigrationRegistry())
	a.runMigration("monitor", a.MonitorDB, a.buildMonitorMigrationRegistry())
}

func (a *App) runMigration(name string, db database.DB, registry *migration.MigrationRegistry) {
	if db == nil {
		zap.L().Error("database is not configured", zap.String("database", name))
		return
	}
	if err := migration.RunMigrations(db.GetDB(), registry); err != nil {
		zap.L().Error("failed to run database migrations",
			zap.String("database", name),
			zap.Error(err),
		)
	}
}

func (a *App) Migrate() error {
	items := []struct {
		db       database.DB
		registry *migration.MigrationRegistry
	}{
		{a.AgentDB, a.buildAgentMigrationRegistry()},
		{a.AppDB, a.buildApplicationMigrationRegistry()},
		{a.ScriptTaskDB, a.buildScriptTaskMigrationRegistry()},
		{a.MonitorDB, a.buildMonitorMigrationRegistry()},
	}
	for _, item := range items {
		if item.db == nil {
			continue
		}
		if err := migration.RunMigrations(item.db.GetDB(), item.registry); err != nil {
			return err
		}
	}
	return nil
}

// RollbackMigration currently targets the primary agent database. Module
// migrations will introduce an explicit logical-database selector in the
// database compatibility stage.
func (a *App) RollbackMigration(version string) error {
	return migration.RollbackMigration(
		a.AgentDB.GetDB(),
		a.buildAgentMigrationRegistry(),
		version,
	)
}

func (a *App) buildAgentMigrationRegistry() *migration.MigrationRegistry {
	registry := migration.NewMigrationRegistry()
	configModule.RegisterMigrations(registry)
	return registry
}

func (a *App) buildApplicationMigrationRegistry() *migration.MigrationRegistry {
	registry := migration.NewMigrationRegistry()
	applicationModule.RegisterMigrations(registry)
	return registry
}

func (a *App) buildScriptTaskMigrationRegistry() *migration.MigrationRegistry {
	registry := migration.NewMigrationRegistry()
	scriptModule.RegisterMigrations(registry)
	return registry
}

func (a *App) buildMonitorMigrationRegistry() *migration.MigrationRegistry {
	registry := migration.NewMigrationRegistry()
	monitorModule.RegisterMigrations(registry)
	return registry
}
