package server

import (
	"squirrel-dev/internal/squ-apiserver/model/application"
	"squirrel-dev/internal/squ-apiserver/model/auth"
	"squirrel-dev/internal/squ-apiserver/model/config"
	"squirrel-dev/internal/squ-apiserver/model/health"
	"squirrel-dev/internal/squ-apiserver/model/migration"
	serverModel "squirrel-dev/internal/squ-apiserver/model/server"

	"go.uber.org/zap"
)

func (s *Server) migrate() {
	registry := migration.NewMigrationRegistry()

	health.RegisterMigrations(registry)
	auth.RegisterMigrations(registry)
	serverModel.RegisterMigrations(registry)
	config.RegisterMigrations(registry)
	application.RegisterMigrations(registry)

	if err := migration.RunMigrations(s.DB.GetDB(), registry); err != nil {
		zap.S().Errorf("迁移失败: %v", err)
	} else {
		zap.S().Infof("迁移成功完成")
	}

}
