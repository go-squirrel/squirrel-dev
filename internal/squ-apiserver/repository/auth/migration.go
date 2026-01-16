package auth

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"用户列表",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&model.User{})
			if err != nil {
				return err
			}
			user := &model.User{
				Username: "test",
				Password: "test",
			}
			return db.Create(user).Error
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("users")
		},
	)
}
