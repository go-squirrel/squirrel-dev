package config

import (
	"path/filepath"
	"testing"
)

func TestLoadLegacyAPIServerConfig(t *testing.T) {
	configPath := filepath.Join("..", "..", "..", "config", "apiserver.yaml")
	value := New(configPath)

	if value.Server.Port != "10700" {
		t.Fatalf("server port = %q, want %q", value.Server.Port, "10700")
	}
	if value.DB.Sqlite.FilePath != "./db/apiserver.db" {
		t.Fatalf("DB path = %q", value.DB.Sqlite.FilePath)
	}
	if value.Auth.Jwt.SigningKey == "" || value.Auth.Jwt.Expired != 1440 {
		t.Fatalf("unexpected auth config: %#v", value.Auth.Jwt)
	}
	if value.MTLS.CAFile != "./certs/ca.crt" {
		t.Fatalf("mTLS CA path = %q", value.MTLS.CAFile)
	}
}
