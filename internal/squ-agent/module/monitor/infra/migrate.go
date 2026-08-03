package infra

import (
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	registry.Register(
		"1.0.0",
		"监控列表",
		func(db *gorm.DB) error {
			for _, model := range []any{&baseMonitorModel{}, &networkMonitorModel{}, &diskIOMonitorModel{}, &diskUsageMonitorModel{}} {
				if err := db.AutoMigrate(model); err != nil {
					return err
				}
			}
			return nil
		},
		func(db *gorm.DB) error {
			// Keep the legacy rollback target exactly as implemented.
			return db.Migrator().DropTable("servers")
		},
	)
}
