package app

import (
	"squirrel-dev/internal/pkg/migration"
	applicationModule "squirrel-dev/internal/squ-apiserver/module/application"
	appstoreModule "squirrel-dev/internal/squ-apiserver/module/appstore"
	authModule "squirrel-dev/internal/squ-apiserver/module/auth"
	configModule "squirrel-dev/internal/squ-apiserver/module/config"
	deploymentModule "squirrel-dev/internal/squ-apiserver/module/deployment"
	scriptModule "squirrel-dev/internal/squ-apiserver/module/script"
	serverModule "squirrel-dev/internal/squ-apiserver/module/server"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// runMigrations preserves the legacy startup behavior: migration failures are
// logged, but they do not prevent the HTTP server from starting.
func (a *App) runMigrations() {
	if a.DB == nil {
		zap.L().Error("database is not configured")
		return
	}
	if err := migration.RunMigrations(a.DB.GetDB(), a.buildMigrationRegistry()); err != nil {
		zap.L().Error("failed to run database migrations", zap.Error(err))
	}
}

func (a *App) Migrate() error {
	return migration.RunMigrations(a.DB.GetDB(), a.buildMigrationRegistry())
}

func (a *App) RollbackMigration(version string) error {
	return migration.RollbackMigration(a.DB.GetDB(), a.buildMigrationRegistry(), version)
}

func (a *App) buildMigrationRegistry() *migration.MigrationRegistry {
	registry := migration.NewMigrationRegistry()
	registry.Register(
		"1.0.0",
		"legacy apiserver schema",
		migrate,
		rollback,
	)
	registry.Register(
		"1.0.1",
		"legacy script execution results",
		scriptModule.MigrateResults,
		scriptModule.RollbackResults,
	)
	return registry
}

func migrate(db *gorm.DB) error {
	if err := authModule.Migrate(db); err != nil {
		return err
	}
	if err := serverModule.Migrate(db); err != nil {
		return err
	}
	if err := configModule.Migrate(db); err != nil {
		return err
	}
	if err := applicationModule.Migrate(db); err != nil {
		return err
	}
	if err := deploymentModule.Migrate(db); err != nil {
		return err
	}
	if err := appstoreModule.Migrate(db); err != nil {
		return err
	}
	return scriptModule.MigrateScripts(db)
}

func rollback(db *gorm.DB) error {
	if err := serverModule.Rollback(db); err != nil {
		return err
	}
	if err := configModule.Rollback(db); err != nil {
		return err
	}
	if err := appstoreModule.Rollback(db); err != nil {
		return err
	}
	if err := applicationModule.Rollback(db); err != nil {
		return err
	}
	if err := deploymentModule.Rollback(db); err != nil {
		return err
	}
	if err := scriptModule.RollbackScripts(db); err != nil {
		return err
	}
	return authModule.Rollback(db)
}
