package infra

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/module/application/domain"
)

func TestLegacyApplicationMigrationAndPersistence(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	registry := migration.NewMigrationRegistry()
	RegisterMigrations(registry)
	if err := migration.RunMigrations(db, registry); err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable("applications") {
		t.Fatal("legacy applications table was not created")
	}

	repository := NewRepository(db)
	value := domain.Application{
		Name: "demo", Status: domain.StatusStopped, DeployID: 99,
		Env: []map[string]string{{"A": "B"}},
	}
	if err := repository.Add(context.Background(), &value); err != nil {
		t.Fatal(err)
	}
	stored, err := repository.GetByDeployID(context.Background(), 99)
	if err != nil {
		t.Fatal(err)
	}
	if stored.ID == 0 || len(stored.Env) != 1 || stored.Env[0]["A"] != "B" {
		t.Fatalf("legacy JSON persistence mismatch: %#v", stored)
	}
}
