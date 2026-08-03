package infra

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestLegacyScriptMigrations(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := MigrateScripts(db); err != nil {
		t.Fatal(err)
	}
	var script scriptModel
	if err := db.Where("name = ?", "test-loop").First(&script).Error; err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(script.Content, "#!") {
		t.Fatalf("embedded seed content changed: %q", script.Content)
	}
	if err := MigrateResults(db); err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable("script_results") {
		t.Fatal("script_results table not created")
	}
}
