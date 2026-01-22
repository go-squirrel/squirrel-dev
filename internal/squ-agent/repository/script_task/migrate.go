package script_task

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册脚本执行任务表初始迁移
	registry.Register(
		"1.0.0",
		"脚本执行任务表",
		// 升级函数
		func(db *gorm.DB) error {
			return db.AutoMigrate(&model.ScriptExecutionTask{})
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("script_execution_tasks")
		},
	)
}
