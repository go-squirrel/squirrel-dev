package application

import (
	"fmt"
	"os"
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
		return "", "", fmt.Errorf("docker-compose 命令不可用")
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

	zap.L().Info("执行 docker-compose",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
		zap.Strings("actions", actions),
	)

	// 切换到工作目录执行命令
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}
	defer os.Chdir(currentDir)

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("切换到工作目录失败: %w", err)
	}

	// 执行命令并获取输出和错误
	stdout, stderr, err := execute.CommandError(command, args...)
	if err != nil {
		zap.L().Error("docker-compose 命令执行失败",
			zap.Strings("actions", actions),
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose 命令执行成功",
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
