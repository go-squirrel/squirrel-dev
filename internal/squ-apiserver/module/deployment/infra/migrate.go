package infra

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error  { return db.AutoMigrate(&deploymentModel{}) }
func Rollback(db *gorm.DB) error { return db.Migrator().DropTable("deployments") }
