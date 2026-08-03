package infra

import (
	"gorm.io/gorm"

	"squirrel-dev/pkg/hash"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&userModel{}); err != nil {
		return err
	}
	password, err := hash.HashPassword("squ123")
	if err != nil {
		return err
	}
	return db.Create(&userModel{Username: "demo", Password: password}).Error
}

func Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable("users")
}
