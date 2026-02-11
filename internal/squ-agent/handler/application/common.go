package application

import (
	"fmt"
	"os"
	"path/filepath"
	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/execute"

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

// getComposePath 获取 docker-compose 文件存储路径
func (a *Application) getComposePath() string {
	if a.Config.Common.ComposePath != "" {
		return a.Config.Common.ComposePath
	}
	return ""
}

// runDockerComposeUp 执行 docker-compose up -d 命令
func runDockerComposeUp(workDir, composeFile string) error {
	return runDockerComposeCommand(workDir, composeFile, "up", "-d")
}

// prepareComposePath 准备 docker-compose 文件目录并返回路径
func (a *Application) prepareComposePath() (string, error) {
	composePath := a.getComposePathOrDefault()

	if err := os.MkdirAll(composePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create compose directory: %w", err)
	}
	return composePath, nil
}

// getComposePathOrDefault 获取 compose 路径，如果为空则返回默认值
func (a *Application) getComposePathOrDefault() string {
	composePath := a.getComposePath()
	if composePath == "" {
		return "."
	}
	return composePath
}

// createOrUpdateComposeFile 创建或更新 docker-compose 文件
func (a *Application) createOrUpdateComposeFile(name, content, composePath string) (string, error) {
	composeFileName := fmt.Sprintf("docker-compose-%s.yml", name)
	composeFilePath := filepath.Join(composePath, composeFileName)

	// 如果文件已存在，先删除（支持重试）
	if _, err := os.Stat(composeFilePath); err == nil {
		zap.L().Info("docker-compose file already exists, deleting it",
			zap.String("path", composeFilePath),
			zap.String("name", name))
		if err := os.Remove(composeFilePath); err != nil {
			return "", fmt.Errorf("failed to delete existing docker-compose file: %w", err)
		}
	}

	if err := os.WriteFile(composeFilePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create docker-compose file: %w", err)
	}

	zap.L().Info("docker-compose file created/updated",
		zap.String("path", composeFilePath),
		zap.String("name", name))

	return composeFilePath, nil
}

// deployApplication 部署应用（创建 compose 文件并启动）
func (a *Application) deployApplication(app *model.Application) (composePath, composeFileName string, err error) {
	// 准备 compose 目录
	composePath, err = a.prepareComposePath()
	if err != nil {
		return "", "", fmt.Errorf("failed to prepare compose path: %w", err)
	}

	// 创建/更新 compose 文件
	_, err = a.createOrUpdateComposeFile(app.Name, app.Content, composePath)
	if err != nil {
		return "", "", err
	}

	composeFileName = fmt.Sprintf("docker-compose-%s.yml", app.Name)

	// 更新数据库状态为 starting
	app.Status = model.AppStatusStarting
	if updateErr := a.Repository.Update(app); updateErr != nil {
		return "", "", fmt.Errorf("failed to update application status: %w", updateErr)
	}

	// 异步启动应用
	go func(appName, composePath, composeFileName string, deployID uint64) {
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
	}(app.Name, composePath, composeFileName, app.DeployID)

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
