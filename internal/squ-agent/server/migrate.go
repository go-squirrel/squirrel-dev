package server

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/repository/application"

	"go.uber.org/zap"
)

func (s *Server) migrate() {
	appRegistry := migration.NewMigrationRegistry()

	application.RegisterMigrations(appRegistry)

	if err := migration.RunMigrations(s.AppDB.GetDB(), appRegistry); err != nil {
		zap.S().Errorf("迁移失败: %v", err)
	} else {
		zap.S().Infof("迁移成功完成")
	}

}
