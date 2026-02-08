package cron

import (
	"encoding/json"
	"fmt"
	"strings"

	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/execute"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

// startApp 启动应用状态检测定时任务
// 每 30 秒检测一次应用容器状态
func (c *Cron) startApp() error {
	_, err := c.Cron.AddFunc("*/5 * * * * *", func() {
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
		shouldUpdate := app.Status != status
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
			}
			reportData := ReportApplicationStatus{
				DeployID: updatedApp.DeployID,
				Status:   status,
			}
			apiServerURL := utils.GenAgentUrl(c.Config.Apiserver.Http.Scheme,
				c.Config.Apiserver.Http.Server,
				0,
				c.Config.Apiserver.Http.BaseUri,
				uriAppReport)
			_, err = c.HTTPClient.Post(apiServerURL, reportData, nil)
			if err != nil {
				zap.L().Error("report apiserver",
					zap.Uint("id", app.ID),
					zap.String("name", app.Name),
					zap.String("old_status", app.Status),
					zap.String("new_status", status),
					zap.Error(err),
				)
			}
			zap.L().Info("应用状态已更新",
				zap.Uint("id", app.ID),
				zap.String("name", app.Name),
				zap.String("old_status", app.Status),
				zap.String("new_status", status),
			)

		}
	}
}

// getContainerStatus 获取指定应用名称的容器状态
func (c *Cron) getContainerStatus(appName string) string {
	// 首先尝试使用 docker compose ls 检查
	composeStatus := c.checkComposeStatus("docker", "compose", appName)
	if composeStatus != "" {
		return composeStatus
	}

	// 如果 docker compose 不可用，尝试 docker-compose ls
	composeStatus = c.checkComposeStatus("docker-compose", "", appName)
	if composeStatus != "" {
		return composeStatus
	}

	// 如果都不可用，返回 unknown
	return model.AppStatusFailed
}

// checkComposeStatus 使用指定的 compose 命令检查应用状态
func (c *Cron) checkComposeStatus(command, composePrefix, appName string) string {
	var args []string
	if composePrefix != "" {
		args = []string{composePrefix, "ls", "--all", "--format", "json"}
	} else {
		args = []string{"ls", "--all", "--format", "json"}
	}

	output, stderr, err := execute.CommandError(command, args...)
	if err != nil {
		zap.L().Warn("检查 compose 状态失败",
			zap.String("command", command),
			zap.String("args", fmt.Sprintf("%v", args)),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return model.AppStatusFailed
	}

	// 解析 JSON 数组输出
	// 输出格式为: [{"Name":"compose","Status":"running(1)","ConfigFiles":"/path/to/docker-compose.yml"}]
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &composeProjects); err != nil {
		zap.L().Debug("解析 compose JSON 输出失败",
			zap.String("output", output),
			zap.Error(err),
		)
		return model.AppStatusFailed
	}

	// 构建 compose 文件名用于匹配
	// compose 文件命名规则: docker-compose-{appName}.yml
	expectedComposeFile := fmt.Sprintf("docker-compose-%s.yml", appName)

	for _, project := range composeProjects {
		// 检查 ConfigFiles 是否包含预期的 compose 文件名
		if strings.Contains(project.ConfigFiles, expectedComposeFile) {
			// 解析状态
			// running(1), running(2) 等表示有容器在运行
			statusLower := strings.ToLower(project.Status)
			if strings.HasPrefix(statusLower, "running") {
				return model.AppStatusRunning
			} else if strings.HasPrefix(statusLower, "exited") {
				return model.AppStatusStopped
			} else if strings.HasPrefix(statusLower, "paused") {
				return model.AppStatusPaused
			}
		}
	}

	return model.AppStatusUndeploy
}
