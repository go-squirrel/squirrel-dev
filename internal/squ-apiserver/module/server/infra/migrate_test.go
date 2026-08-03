package infra

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestLegacyServerMigration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable("servers") {
		t.Fatal("legacy servers table was not created")
	}
	var server serverModel
	if err := db.First(&server).Error; err != nil {
		t.Fatal(err)
	}
	if server.UUID == "" || server.Hostname == "" || server.IPAddress != "127.0.0.1" {
		t.Fatalf("legacy seed mismatch: %#v", server)
	}
}
