package config

import (
	"path/filepath"
	"testing"
)

func TestLoadLegacyAgentConfig(t *testing.T) {
	configPath := filepath.Join("..", "..", "..", "config", "agent.yaml")
	value := New(configPath)

	if value.Server.Port != "10750" {
		t.Fatalf("server port = %q, want %q", value.Server.Port, "10750")
	}
	if value.Cache.Type != "memory" {
		t.Fatalf("cache type = %q, want %q", value.Cache.Type, "memory")
	}
	if value.DB.Sqlite.AgentFilePath != "./db/agent.db" {
		t.Fatalf("agent DB path = %q", value.DB.Sqlite.AgentFilePath)
	}
	if value.DB.Sqlite.AppFilePath != "./db/agent-app.db" {
		t.Fatalf("application DB path = %q", value.DB.Sqlite.AppFilePath)
	}
	if value.DB.Sqlite.MonitorFilePath != "./db/agent-monitor.db" {
		t.Fatalf("monitor DB path = %q", value.DB.Sqlite.MonitorFilePath)
	}
	if value.DB.Sqlite.ScriptTaskFilePath != "./db/agent-script-task.db" {
		t.Fatalf("script task DB path = %q", value.DB.Sqlite.ScriptTaskFilePath)
	}
}
