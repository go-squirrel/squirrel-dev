package infra

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/module/script/domain"
)

func TestLegacyScriptTaskMigrationAndQueries(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	registry := migration.NewMigrationRegistry()
	RegisterMigrations(registry)
	if err := migration.RunMigrations(db, registry); err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable("script_execution_tasks") {
		t.Fatal("legacy script_execution_tasks table was not created")
	}

	repository := NewRepository(db)
	task := domain.Task{ScriptID: 2, TaskID: 3, Name: "demo", Status: "success"}
	if err := repository.Add(context.Background(), &task); err != nil {
		t.Fatal(err)
	}
	values, err := repository.GetUnreportedTasks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 1 || values[0].TaskID != 3 {
		t.Fatalf("unexpected unreported tasks: %#v", values)
	}
	if err := repository.MarkAsReported(context.Background(), task.ID); err != nil {
		t.Fatal(err)
	}
	values, err = repository.GetUnreportedTasks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 0 {
		t.Fatalf("reported task still selected: %#v", values)
	}
}
