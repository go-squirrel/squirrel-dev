package monitor

import (
	"squirrel-dev/internal/squ-agent/model/migration"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"监控列表",
		// 升级函数
		func(db *gorm.DB) error {
			return db.AutoMigrate(&Monitor{})
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("servers")
		},
	)
}
