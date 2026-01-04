package server

import (
	"squirrel-dev/internal/squ-apiserver/model/migration"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"初始化服务器",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&Server{})
			if err != nil {
				return err
			}
			return err
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("servers")
		},
	)
}
