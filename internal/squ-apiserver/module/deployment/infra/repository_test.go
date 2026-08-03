package infra

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/deployment/domain"
)

func TestLegacyDeploymentPersistence(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	repository := NewRepository(db)
	value := domain.Deployment{ServerID: 2, ApplicationID: 3, DeployID: 99, Env: []map[string]string{{"A": "B"}}}
	if err := repository.Add(context.Background(), &value); err != nil {
		t.Fatal(err)
	}
	stored, err := repository.GetByDeployID(context.Background(), 99)
	if err != nil {
		t.Fatal(err)
	}
	if stored.ID == 0 || stored.Env[0]["A"] != "B" {
		t.Fatalf("stored = %#v", stored)
	}
}
