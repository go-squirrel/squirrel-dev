package deployment

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册应用服务器关联表初始迁移
	registry.Register(
		"1.0.0",
		"应用服务器关联表",
		// 升级函数
		func(db *gorm.DB) error {
			return db.AutoMigrate(&model.Deployment{})
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("deployments")
		},
	)
}
