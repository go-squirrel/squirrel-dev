package server

import (
	"squirrel-dev/internal/apiserver/model/health"
	"squirrel-dev/internal/apiserver/model/migration"

	"go.uber.org/zap"
)

func (s *Server) migrate() {
	registry := migration.NewMigrationRegistry()

	health.RegisterMigrations(registry)

	if err := migration.RunMigrations(s.DB.GetDB(), registry); err != nil {
		zap.S().Errorf("迁移失败: %v", err)
	} else {
		zap.S().Infof("迁移成功完成")
	}

}
