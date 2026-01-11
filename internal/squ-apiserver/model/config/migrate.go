package config

import (
	"squirrel-dev/internal/squ-apiserver/model/migration"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"配置列表",
		// 升级函数
		func(db *gorm.DB) error {
			db.AutoMigrate(&Config{})
			c := []Config{
				{Key: "registry", Value: "docker.io"},
				{Key: "registry_username", Value: ""},
				{Key: "registry_password", Value: ""},
			}
			return db.Create(c).Error
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("servers")
		},
	)
}
