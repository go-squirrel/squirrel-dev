package infra

import (
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	registry.Register(
		"1.0.0",
		"配置列表",
		func(db *gorm.DB) error {
			if err := db.AutoMigrate(&configModel{}); err != nil {
				return err
			}
			return db.Create(&[]configModel{
				{Key: "monitor_interval", Value: "300"},
				{Key: "monitor_expired", Value: "604800"},
			}).Error
		},
		func(db *gorm.DB) error {
			// Keep the legacy rollback target exactly as implemented.
			return db.Migrator().DropTable("servers")
		},
	)
}
