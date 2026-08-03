package infra

import (
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	registry.Register(
		"1.0.0",
		"脚本执行任务表",
		func(db *gorm.DB) error { return db.AutoMigrate(&taskModel{}) },
		func(db *gorm.DB) error { return db.Migrator().DropTable("script_execution_tasks") },
	)
}
