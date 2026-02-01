package auth

import (
	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/hash"

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
			hashedPassword, err := hash.HashPassword("squ123")
			if err != nil {
				return err
			}
			user := &model.User{
				Username: "demo",
				Password: hashedPassword,
			}
			return db.Create(user).Error
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("users")
		},
	)
}
