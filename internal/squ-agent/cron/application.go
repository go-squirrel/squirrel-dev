package cron

import (
	"fmt"
	"strings"

	"squirrel-dev/pkg/execute"

	"go.uber.org/zap"
)

// startApp 启动应用状态检测定时任务
// 每 30 秒检测一次应用容器状态
func (c *Cron) startApp() error {
	_, err := c.Cron.AddFunc("*/30 * * * * *", func() {
		c.checkApplicationStatus()
	})
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

// checkApplicationStatus 检查所有应用的容器状态
func (c *Cron) checkApplicationStatus() {
	// 获取所有应用
	applications, err := c.AppRepository.List()
	if err != nil {
		zap.L().Error("获取应用列表失败", zap.Error(err))
		return
	}

	// 遍历每个应用，检查容器状态
	for _, app := range applications {
		status := c.getContainerStatus(app.Name)

		// 对于 "starting" 状态，无论检测结果如何，都更新数据库
		// 对于 "failed" 状态，也需要重新检测
		// 对于其他稳定状态，只有状态发生变化时才更新
		shouldUpdate := app.Status != status || app.Status == "starting" || app.Status == "failed"
		if shouldUpdate {
			updatedApp := app
			updatedApp.Status = status
			err := c.AppRepository.Update(&updatedApp)
			if err != nil {
				zap.L().Error("更新应用状态失败",
					zap.Uint("id", app.ID),
					zap.String("name", app.Name),
					zap.String("old_status", app.Status),
					zap.String("new_status", status),
					zap.Error(err),
				)
			} else {
				zap.L().Info("应用状态已更新",
					zap.Uint("id", app.ID),
					zap.String("name", app.Name),
					zap.String("old_status", app.Status),
					zap.String("new_status", status),
				)
			}
		}
	}
}

// getContainerStatus 获取指定应用名称的容器状态
func (c *Cron) getContainerStatus(appName string) string {
	// 使用 docker ps 命令检查容器状态
	// 查找名称包含 appName 的容器
	output, stderr, err := execute.CommandError("docker", "ps", "--format", "{{.Names}}:{{.Status}}", "--filter", fmt.Sprintf("name=%s", appName))

	if err != nil {
		zap.L().Warn("检查容器状态失败",
			zap.String("app_name", appName),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return "unknown"
	}

	// 解析输出，查找匹配的容器
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 输出格式为：container_name:status
		// 例如：my-app-1:Up 2 hours
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			containerName := parts[0]
			status := parts[1]

			// 检查容器名称是否匹配
			if strings.Contains(containerName, appName) {
				// 判断容器状态
				if strings.HasPrefix(status, "Up") {
					return "running"
				} else if strings.HasPrefix(status, "Exited") {
					return "stopped"
				} else if strings.HasPrefix(status, "Restarting") {
					return "restarting"
				}
			}
		}
	}

	// 没有找到运行中的容器，检查是否已停止
	output, stderr, err = execute.CommandError("docker", "ps", "-a", "--format", "{{.Names}}:{{.Status}}", "--filter", fmt.Sprintf("name=%s", appName))

	if err != nil {
		return "unknown"
	}

	lines = strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			containerName := parts[0]
			status := parts[1]

			if strings.Contains(containerName, appName) {
				if strings.HasPrefix(status, "Exited") {
					return "stopped"
				}
			}
		}
	}

	// 没有找到容器
	return "not_deployed"
}
