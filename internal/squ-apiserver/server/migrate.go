package server

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-apiserver/repository/app_store"
	"squirrel-dev/internal/squ-apiserver/repository/application"
	"squirrel-dev/internal/squ-apiserver/repository/auth"
	"squirrel-dev/internal/squ-apiserver/repository/config"
	"squirrel-dev/internal/squ-apiserver/repository/deployment"
	"squirrel-dev/internal/squ-apiserver/repository/health"
	"squirrel-dev/internal/squ-apiserver/repository/script"
	serverModel "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
)

func (s *Server) migrate() {
	registry := migration.NewMigrationRegistry()

	health.RegisterMigrations(registry)
	auth.RegisterMigrations(registry)
	serverModel.RegisterMigrations(registry)
	config.RegisterMigrations(registry)
	application.RegisterMigrations(registry)
	deployment.RegisterMigrations(registry)
	app_store.RegisterMigrations(registry)
	script.RegisterMigrations(registry)

	if err := migration.RunMigrations(s.DB.GetDB(), registry); err != nil {
		zap.S().Errorf("迁移失败: %v", err)
	} else {
		zap.S().Infof("迁移成功完成")
	}

}
