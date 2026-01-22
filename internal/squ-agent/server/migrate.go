package server

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/repository/application"
	"squirrel-dev/internal/squ-agent/repository/script_task"

	"go.uber.org/zap"
)

func (s *Server) migrate() {
	// 迁移 AppDB
	appRegistry := migration.NewMigrationRegistry()
	application.RegisterMigrations(appRegistry)

	if err := migration.RunMigrations(s.AppDB.GetDB(), appRegistry); err != nil {
		zap.S().Errorf("AppDB 迁移失败: %v", err)
	} else {
		zap.S().Infof("AppDB 迁移成功完成")
	}

	// 迁移 ScriptTaskDB
	scriptTaskRegistry := migration.NewMigrationRegistry()
	script_task.RegisterMigrations(scriptTaskRegistry)

	if err := migration.RunMigrations(s.ScriptTaskDB.GetDB(), scriptTaskRegistry); err != nil {
		zap.S().Errorf("ScriptTaskDB 迁移失败: %v", err)
	} else {
		zap.S().Infof("ScriptTaskDB 迁移成功完成")
	}

}
