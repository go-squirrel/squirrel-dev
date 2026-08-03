package infra

import (
	"fmt"
	"os"

	"squirrel-dev/pkg/execute"
)

type ShellExecutor struct{}

func NewShellExecutor() *ShellExecutor { return &ShellExecutor{} }

func (ShellExecutor) Execute(taskID uint, content string) (string, error) {
	path := fmt.Sprintf("/tmp/script_%d.sh", taskID)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		return "", err
	}
	defer os.Remove(path)
	return execute.Command("bash", path)
}
