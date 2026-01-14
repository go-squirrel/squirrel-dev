package app_store

import (
	"squirrel-dev/internal/squ-apiserver/model/migration"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册应用商店表迁移
	registry.Register(
		"1.0.0",
		"应用商店",
		// 升级函数
		func(db *gorm.DB) error {
			return db.AutoMigrate(&AppStore{})
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("app_stores")
		},
	)
}
