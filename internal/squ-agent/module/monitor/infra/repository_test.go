package infra

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-agent/module/monitor/domain"
)

func TestLegacyMonitorMigrationAndRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	registry := migration.NewMigrationRegistry()
	RegisterMigrations(registry)
	if err := migration.RunMigrations(db, registry); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"base_monitors", "network_monitors", "disk_io_monitors", "disk_usage_monitors"} {
		if !db.Migrator().HasTable(table) {
			t.Fatalf("missing legacy table %q", table)
		}
	}

	repository := NewRepository(db)
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	old := domain.BaseMonitor{CPUUsage: 10, CollectTime: now.Add(-2 * time.Hour)}
	recent := domain.BaseMonitor{CPUUsage: 20, CollectTime: now.Add(-time.Minute)}
	if err := repository.CreateBaseMonitor(context.Background(), &recent); err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateBaseMonitor(context.Background(), &old); err != nil {
		t.Fatal(err)
	}

	values, err := repository.BaseByTimeRange(context.Background(), now.Add(-3*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 2 || values[0].CPUUsage != 10 || values[1].CPUUsage != 20 {
		t.Fatalf("values are not ordered by collect_time: %#v", values)
	}

	if err := repository.DeleteBeforeTime(context.Background(), now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	values, err = repository.BaseByTimeRange(context.Background(), now.Add(-3*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 1 || values[0].CPUUsage != 20 {
		t.Fatalf("legacy soft delete mismatch: %#v", values)
	}
}
