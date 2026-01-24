package config

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"配置列表",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&model.Config{})
			if err != nil {
				return err
			}
			c := []model.Config{
				{
					Key:   "monitor_interval",
					Value: "300",
				},
				{
					Key:   "monitor_expired",
					Value: "604800",
				},
			}

			return db.Create(&c).Error
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("servers")
		},
	)
}
