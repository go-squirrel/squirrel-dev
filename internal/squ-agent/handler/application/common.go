package application

import (
	"fmt"
	"os"
	"squirrel-dev/pkg/execute"

	"go.uber.org/zap"
)

// runDockerComposeStart 执行 docker-compose start 命令
func runDockerComposeStart(workDir, composeFile string) error {
	var command string
	var args []string

	// 尝试使用 docker-compose 命令
	if checkCommandAvailable("docker-compose") {
		command = "docker-compose"
		args = []string{"-f", composeFile, "start"}
	} else if checkCommandAvailable("docker") {
		// 使用 docker compose 插件
		command = "docker"
		args = []string{"compose", "-f", composeFile, "start"}
	} else {
		return fmt.Errorf("docker-compose 命令不可用")
	}

	zap.L().Info("执行 docker-compose start",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
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
		zap.L().Error("docker-compose start 命令执行失败",
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose start 执行成功",
		zap.String("output", stdout),
	)

	return nil
}

// runDockerComposeStop 执行 docker-compose stop 命令
func runDockerComposeStop(workDir, composeFile string) error {
	var command string
	var args []string

	// 尝试使用 docker-compose 命令
	if checkCommandAvailable("docker-compose") {
		command = "docker-compose"
		args = []string{"-f", composeFile, "stop"}
	} else if checkCommandAvailable("docker") {
		// 使用 docker compose 插件
		command = "docker"
		args = []string{"compose", "-f", composeFile, "stop"}
	} else {
		return fmt.Errorf("docker-compose 命令不可用")
	}

	zap.L().Info("执行 docker-compose stop",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
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
		zap.L().Error("docker-compose stop 命令执行失败",
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose stop 执行成功",
		zap.String("output", stdout),
	)

	return nil
}

// checkCommandAvailable 检查命令是否在 PATH 中可用
func checkCommandAvailable(command string) bool {
	// 尝试执行命令的 --version 参数来检查命令是否可用
	_, err := execute.Command(command, "--version")
	return err == nil
}

// runDockerComposeUp 执行 docker-compose up -d 命令
func runDockerComposeUp(workDir, composeFile string) error {
	var command string
	var args []string

	// 尝试使用 docker-compose 命令
	if checkCommandAvailable("docker-compose") {
		command = "docker-compose"
		args = []string{"-f", composeFile, "up", "-d"}
	} else if checkCommandAvailable("docker") {
		// 使用 docker compose 插件
		command = "docker"
		args = []string{"compose", "-f", composeFile, "up", "-d"}
	} else {
		return fmt.Errorf("docker-compose 命令不可用")
	}

	zap.L().Info("执行 docker-compose up -d",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
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
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose 启动成功",
		zap.String("output", stdout),
	)

	return nil
}

// checkDockerComposeAvailable 检测 docker-compose 命令是否可用
func checkDockerComposeAvailable() bool {
	// 检查 docker-compose 命令
	_, err := execute.Command("docker-compose", "--version")
	if err == nil {
		return true
	}

	// 检查 docker compose 插件
	_, err = execute.Command("docker", "compose", "version")
	if err == nil {
		return true
	}

	return false
}

// getComposePath 获取 docker-compose 文件存储路径
func (a *Application) getComposePath() string {
	if a.Config.Common.ComposePath != "" {
		return a.Config.Common.ComposePath
	}
	return ""
}
