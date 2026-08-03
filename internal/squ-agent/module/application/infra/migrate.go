package infra

import (
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	registry.Register(
		"1.0.0",
		"应用列表",
		func(db *gorm.DB) error { return db.AutoMigrate(&applicationModel{}) },
		func(db *gorm.DB) error {
			// Keep the legacy rollback target exactly as implemented.
			return db.Migrator().DropTable("servers")
		},
	)
}
