package server

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/repository/application"
	"squirrel-dev/internal/squ-agent/repository/config"
	"squirrel-dev/internal/squ-agent/repository/monitor"
	"squirrel-dev/internal/squ-agent/repository/script_task"

	"go.uber.org/zap"
)

func (s *Server) migrate() {
	// 通用内容
	configRegistry := migration.NewMigrationRegistry()
	config.RegisterMigrations(configRegistry)

	if err := migration.RunMigrations(s.AgentDB.GetDB(), configRegistry); err != nil {
		zap.L().Error("Failed to migrate config", zap.Error(err))
	} else {
		zap.L().Info("Config migration completed successfully")
	}

	// 迁移 AppDB
	appRegistry := migration.NewMigrationRegistry()
	application.RegisterMigrations(appRegistry)

	if err := migration.RunMigrations(s.AppDB.GetDB(), appRegistry); err != nil {
		zap.L().Error("Failed to migrate AppDB", zap.Error(err))
	} else {
		zap.L().Info("AppDB migration completed successfully")
	}

	// 迁移 ScriptTaskDB
	scriptTaskRegistry := migration.NewMigrationRegistry()
	script_task.RegisterMigrations(scriptTaskRegistry)

	if err := migration.RunMigrations(s.ScriptTaskDB.GetDB(), scriptTaskRegistry); err != nil {
		zap.L().Error("Failed to migrate ScriptTaskDB", zap.Error(err))
	} else {
		zap.L().Info("ScriptTaskDB migration completed successfully")
	}
	// 迁移 MonitorDB
	monitorRegistry := migration.NewMigrationRegistry()
	monitor.RegisterMigrations(monitorRegistry)

	if err := migration.RunMigrations(s.MonitorDB.GetDB(), monitorRegistry); err != nil {
		zap.L().Error("Failed to migrate MonitorDB", zap.Error(err))
	} else {
		zap.L().Info("MonitorDB migration completed successfully")
	}
}
