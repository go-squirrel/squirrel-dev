package app

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLegacyCLIContract(t *testing.T) {
	command := NewServerCommand()

	flag := command.Flags().Lookup("config")
	if flag == nil || flag.DefValue != "config/squctl.yaml" {
		t.Fatalf("config default = %#v", flag)
	}

	assertCommand(t, command, "version")
	certs := assertCommand(t, command, "certs")
	assertNoCommand(t, command, "migrate")
	assertNoCommand(t, command, "rollback")

	expectedDefaults := map[string]string{
		"output":       "./certs",
		"ca-cn":        "squirrel-ca",
		"server-cn":    "squirrel-apiserver",
		"server-hosts": "[127.0.0.1,localhost]",
		"client-cn":    "squirrel-agent",
		"expiry":       "87600h0m0s",
		"key-size":     "2048",
	}
	for name, expected := range expectedDefaults {
		flag := certs.Flags().Lookup(name)
		if flag == nil || flag.DefValue != expected {
			t.Fatalf("%s default = %#v, want %q", name, flag, expected)
		}
	}
}

func assertCommand(t *testing.T, command *cobra.Command, name string) *cobra.Command {
	t.Helper()
	for _, child := range command.Commands() {
		if child.Name() == name {
			return child
		}
	}
	t.Fatalf("command %q is missing", name)
	return nil
}

func assertNoCommand(t *testing.T, command *cobra.Command, name string) {
	t.Helper()
	for _, child := range command.Commands() {
		if child.Name() == name {
			t.Fatalf("command %q must not be registered", name)
		}
	}
}
