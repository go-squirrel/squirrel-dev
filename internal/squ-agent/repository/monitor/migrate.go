package monitor

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"监控列表",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&model.BaseMonitor{})
			if err != nil {
				return err
			}
			err = db.AutoMigrate(&model.NetworkMonitor{})
			if err != nil {
				return err
			}
			err = db.AutoMigrate(&model.DiskIOMonitor{})
			if err != nil {
				return err
			}
			err = db.AutoMigrate(&model.DiskUsageMonitor{})
			return err
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("servers")
		},
	)
}
