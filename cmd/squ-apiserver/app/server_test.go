package app

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLegacyCLIContract(t *testing.T) {
	command := NewServerCommand()

	flag := command.Flags().Lookup("config")
	if flag == nil || flag.DefValue != "config/apiserver.yaml" {
		t.Fatalf("config default = %#v", flag)
	}

	assertCommand(t, command, "version")
	assertCommand(t, command, "migrate")
	assertCommand(t, command, "rollback")
}

func assertCommand(t *testing.T, command *cobra.Command, name string) {
	t.Helper()
	for _, child := range command.Commands() {
		if child.Name() == name {
			return
		}
	}
	t.Fatalf("command %q is missing", name)
}
