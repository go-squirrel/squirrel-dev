package infra

import (
	"os"

	"gorm.io/gorm"

	"squirrel-dev/pkg/utils"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&serverModel{}); err != nil {
		return err
	}
	hostname, _ := os.Hostname()
	return db.Create(&serverModel{
		UUID: utils.GenerateServerUUID(hostname), Hostname: hostname, IPAddress: "127.0.0.1",
	}).Error
}

func Rollback(db *gorm.DB) error { return db.Migrator().DropTable("servers") }
