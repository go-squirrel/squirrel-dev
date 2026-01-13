package auth

import (
	"squirrel-dev/internal/squ-apiserver/model/migration"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册配置表初始迁移
	registry.Register(
		"1.0.0",
		"用户列表",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&User{})
			if err != nil {
				return err
			}
			user := &User{
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
