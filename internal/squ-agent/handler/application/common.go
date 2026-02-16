package application

import (
	"fmt"
	"os"
	"path/filepath"
	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/execute"
	"strings"

	"go.uber.org/zap"
)

// getDockerComposeCommandType 获取可用的 docker-compose 命令类型
// 返回使用的命令和 compose 前缀
func getDockerComposeCommandType() (command string, composePrefix string, err error) {
	if checkCommandAvailable("docker-compose") {
		return "docker-compose", "", nil
	} else if checkCommandAvailable("docker") {
		return "docker", "compose", nil
	} else {
		return "", "", fmt.Errorf("docker-compose command not available")
	}
}

// runDockerComposeCommand 通用的 docker-compose 命令执行函数
func runDockerComposeCommand(workDir, composeFile string, actions ...string) error {
	command, composePrefix, err := getDockerComposeCommandType()
	if err != nil {
		return err
	}

	var args []string
	if composePrefix != "" {
		args = append([]string{composePrefix, "-f", composeFile}, actions...)
	} else {
		args = append([]string{"-f", composeFile}, actions...)
	}

	zap.L().Info("Executing docker-compose",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
		zap.Strings("actions", actions),
	)

	// 切换到工作目录执行命令
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(currentDir)

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("failed to change to work directory: %w", err)
	}

	// 执行命令并获取输出和错误
	stdout, stderr, err := execute.CommandError(command, args...)
	if err != nil {
		zap.L().Error("docker-compose command execution failed",
			zap.Strings("actions", actions),
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose command executed successfully",
		zap.Strings("actions", actions),
		zap.String("output", stdout),
	)

	return nil
}

// runDockerComposeStart 执行 docker-compose start 命令
func runDockerComposeStart(workDir, composeFile string) error {
	return runDockerComposeCommand(workDir, composeFile, "start")
}

// runDockerComposeStop 执行 docker-compose stop 命令
func runDockerComposeStop(workDir, composeFile string) error {
	return runDockerComposeCommand(workDir, composeFile, "stop")
}

// checkCommandAvailable 检查命令是否在 PATH 中可用
func checkCommandAvailable(command string) bool {
	// 尝试执行命令的 --version 参数来检查命令是否可用
	_, err := execute.Command(command, "--version")
	return err == nil
}

// checkDockerComposeAvailable 检测 docker-compose 命令是否可用
func checkDockerComposeAvailable() bool {
	_, _, err := getDockerComposeCommandType()
	return err == nil
}

// runDockerComposeUp 执行 docker-compose up -d 命令
func runDockerComposeUp(workDir, composeFile string) error {
	return runDockerComposeCommand(workDir, composeFile, "up", "-d")
}

func runDockerComposeDown(workDir, composeFile string) error {
	return runDockerComposeCommand(workDir, composeFile, "down")
}

// checkDockerInstalled 检测 Docker 是否已安装
func checkDockerInstalled() bool {
	_, err := execute.Command("docker", "--version")
	return err == nil
}

// prepareComposePath 准备 docker-compose 文件目录并返回路径
// 目录结构: composePath/deployID/
func (a *Application) prepareComposePath(deployID uint64) (string, error) {
	composePath := a.getComposePathOrDefault()
	appComposePath := filepath.Join(composePath, fmt.Sprintf("%d", deployID))

	if err := os.MkdirAll(appComposePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create compose directory: %w", err)
	}
	return appComposePath, nil
}

// getComposePathOrDefault 获取 compose 路径，如果为空则返回默认值
func (a *Application) getComposePathOrDefault() string {
	composePath := a.Config.Common.ComposePath
	if composePath == "" {
		return "."
	}
	return composePath
}

// createOrUpdateComposeFile 创建或更新 docker-compose 文件
// 文件名固定为 docker-compose.yml
func (a *Application) createOrUpdateComposeFile(content, composePath string) (string, error) {
	composeFileName := "docker-compose.yml"
	composeFilePath := filepath.Join(composePath, composeFileName)

	// 如果文件已存在，先删除（支持重试）
	if _, err := os.Stat(composeFilePath); err == nil {
		zap.L().Info("docker-compose file already exists, deleting it",
			zap.String("path", composeFilePath))
		if err := os.Remove(composeFilePath); err != nil {
			return "", fmt.Errorf("failed to delete existing docker-compose file: %w", err)
		}
	}

	if err := os.WriteFile(composeFilePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create docker-compose file: %w", err)
	}

	zap.L().Info("docker-compose file created/updated",
		zap.String("path", composeFilePath))

	return composeFilePath, nil
}

// createOrUpdateEnvFile 创建或更新 .env 文件
// env 格式: []map[string]string，每个 map 只有一个 key-value 对
func (a *Application) createOrUpdateEnvFile(env []map[string]string, composePath string) error {
	if len(env) == 0 {
		return nil
	}

	envFilePath := filepath.Join(composePath, ".env")

	var lines []string
	for _, item := range env {
		for key, value := range item {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
		}
	}

	content := strings.Join(lines, "\n")
	if len(content) > 0 {
		content += "\n"
	}

	if err := os.WriteFile(envFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	zap.L().Info(".env file created/updated",
		zap.String("path", envFilePath),
		zap.Int("env_count", len(env)))

	return nil
}

// prepareComposeFiles 准备 compose 文件（不涉及数据库操作）
// 返回 compose 目录路径和文件名，目录结构: composePath/deployID/docker-compose.yml
func (a *Application) prepareComposeFiles(app *model.Application) (composePath, composeFileName string, err error) {
	// 准备 compose 目录 (composePath/deployID/)
	composePath, err = a.prepareComposePath(app.DeployID)
	if err != nil {
		return "", "", fmt.Errorf("failed to prepare compose path: %w", err)
	}

	// 创建/更新 compose 文件
	_, err = a.createOrUpdateComposeFile(app.Content, composePath)
	if err != nil {
		return "", "", err
	}

	// 创建/更新 .env 文件
	if err := a.createOrUpdateEnvFile(app.Env, composePath); err != nil {
		return "", "", fmt.Errorf("failed to create .env file: %w", err)
	}

	// 文件名固定为 docker-compose.yml
	composeFileName = "docker-compose.yml"
	return composePath, composeFileName, nil
}

// startDockerComposeAsync 异步启动 docker-compose
func (a *Application) startDockerComposeAsync(appName, composePath, composeFileName string, deployID uint64) {
	go func() {
		zap.L().Info("Starting application asynchronously", zap.String("name", appName), zap.Uint64("deploy_id", deployID))

		if err := runDockerComposeUp(composePath, composeFileName); err != nil {
			zap.L().Error("Failed to start docker-compose",
				zap.String("path", composePath),
				zap.String("file", composeFileName),
				zap.Uint64("deploy_id", deployID),
				zap.Error(err))

			a.updateApplicationStatusToFailed(deployID)
			return
		}

		zap.L().Info("docker-compose up command executed successfully",
			zap.String("name", appName),
			zap.Uint64("deploy_id", deployID))
	}()
}

// deployApplication 部署应用（创建 compose 文件并启动）
func (a *Application) deployApplication(app *model.Application) (composePath, composeFileName string, err error) {
	composePath, composeFileName, err = a.prepareComposeFiles(app)
	if err != nil {
		return "", "", err
	}

	// 更新数据库状态为 starting
	app.Status = model.AppStatusStarting
	if updateErr := a.Repository.Update(app); updateErr != nil {
		return "", "", fmt.Errorf("failed to update application status: %w", updateErr)
	}

	// 异步启动应用
	a.startDockerComposeAsync(app.Name, composePath, composeFileName, app.DeployID)

	return composePath, composeFileName, nil
}

// updateApplicationStatusToFailed 更新应用状态为失败
func (a *Application) updateApplicationStatusToFailed(deployID uint64) {
	apps, err := a.Repository.List()
	if err != nil {
		zap.L().Error("Failed to list applications when updating to failed status", zap.Error(err))
		return
	}
	for i := range apps {
		if apps[i].DeployID == deployID {
			apps[i].Status = model.AppStatusFailed
			if updateErr := a.Repository.Update(&apps[i]); updateErr != nil {
				zap.L().Error("Failed to update application status to failed",
					zap.Uint64("deploy_id", deployID),
					zap.Error(updateErr))
			}
			break
		}
	}
}
