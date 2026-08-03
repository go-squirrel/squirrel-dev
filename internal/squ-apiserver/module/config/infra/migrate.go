package infra

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	// Keep the old ignored AutoMigrate error.
	_ = db.AutoMigrate(&configModel{})
	return db.Create([]configModel{
		{Key: "registry", Value: "docker.io"},
		{Key: "registry_username", Value: ""},
		{Key: "registry_password", Value: ""},
	}).Error
}

func Rollback(db *gorm.DB) error { return db.Migrator().DropTable("servers") }
