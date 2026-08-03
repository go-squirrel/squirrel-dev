package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"squirrel-dev/pkg/execute"
)

type ComposeRuntime struct {
	basePath string
}

func NewComposeRuntime(basePath string) *ComposeRuntime {
	if basePath == "" {
		basePath = "."
	}
	return &ComposeRuntime{basePath: basePath}
}

func (r *ComposeRuntime) DockerInstalled() bool {
	_, err := execute.Command("docker", "--version")
	return err == nil
}

func (r *ComposeRuntime) ComposeAvailable() bool {
	_, _, err := composeCommand()
	return err == nil
}

func (r *ComposeRuntime) Prepare(deployID uint64, content string, env []map[string]string) (string, string, error) {
	path := r.Path(deployID)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create compose directory: %w", err)
	}
	composeFile := filepath.Join(path, "docker-compose.yml")
	if _, err := os.Stat(composeFile); err == nil {
		if err := os.Remove(composeFile); err != nil {
			return "", "", fmt.Errorf("failed to delete existing docker-compose file: %w", err)
		}
	}
	if err := os.WriteFile(composeFile, []byte(content), 0644); err != nil {
		return "", "", fmt.Errorf("failed to create docker-compose file: %w", err)
	}
	if len(env) > 0 {
		var lines []string
		for _, item := range env {
			for key, value := range item {
				lines = append(lines, fmt.Sprintf("%s=%s", key, value))
			}
		}
		envContent := strings.Join(lines, "\n")
		if envContent != "" {
			envContent += "\n"
		}
		if err := os.WriteFile(filepath.Join(path, ".env"), []byte(envContent), 0644); err != nil {
			return "", "", fmt.Errorf("failed to create .env file: %w", err)
		}
	}
	return path, "docker-compose.yml", nil
}

func (r *ComposeRuntime) ComposeFileExists(deployID uint64) bool {
	_, err := os.Stat(filepath.Join(r.Path(deployID), "docker-compose.yml"))
	return !os.IsNotExist(err)
}

func (r *ComposeRuntime) Path(deployID uint64) string {
	return filepath.Join(r.basePath, fmt.Sprintf("%d", deployID))
}

func (r *ComposeRuntime) Up(path, file string) error    { return runCompose(path, file, "up", "-d") }
func (r *ComposeRuntime) Start(path, file string) error { return runCompose(path, file, "start") }
func (r *ComposeRuntime) Stop(path, file string) error  { return runCompose(path, file, "stop") }

func composeCommand() (string, string, error) {
	if _, err := execute.Command("docker-compose", "--version"); err == nil {
		return "docker-compose", "", nil
	}
	if _, err := execute.Command("docker", "--version"); err == nil {
		return "docker", "compose", nil
	}
	return "", "", fmt.Errorf("docker-compose command not available")
}

func runCompose(workDir, composeFile string, actions ...string) error {
	command, prefix, err := composeCommand()
	if err != nil {
		return err
	}
	var args []string
	if prefix == "" {
		args = append([]string{"-f", composeFile}, actions...)
	} else {
		args = append([]string{prefix, "-f", composeFile}, actions...)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(currentDir)
	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("failed to change to work directory: %w", err)
	}
	_, _, err = execute.CommandError(command, args...)
	return err
}
