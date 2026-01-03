package health

import (
	"squirrel-dev/internal/model/migration"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"应用列表",
		// 升级函数
		func(db *gorm.DB) error {
			return db.AutoMigrate(&Health{})
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("servers")
		},
	)
}

