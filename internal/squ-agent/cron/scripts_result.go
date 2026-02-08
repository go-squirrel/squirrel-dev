package cron

import (
	"encoding/json"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

// startScriptResultReporter 启动脚本执行结果上报定时任务
// 每 5 秒上报一次未上报的执行结果
func (c *Cron) startScriptResultReporter() error {
	_, err := c.Cron.AddFunc("*/5 * * * * *", func() {
		c.reportScriptResults()
	})
	if err != nil {
		zap.L().Error("failed to start cron task for script result reporter",
			zap.String("cron", "script_result_reporter"),
			zap.Error(err))
	}
	return err
}

// reportScriptResults 上报脚本执行结果到 API Server
func (c *Cron) reportScriptResults() {
	// 获取未上报的任务
	tasks, err := c.ScriptTaskRepo.GetUnreportedTasks()
	if err != nil {
		zap.L().Error("failed to get unreported script tasks",
			zap.String("cron", "script_result_reporter"),
			zap.Error(err))
		return
	}

	if len(tasks) == 0 {
		return
	}

	zap.L().Info("started reporting script execution results",
		zap.String("cron", "script_result_reporter"),
		zap.Int("count", len(tasks)))
	// TODO: 从配置中获取 API Server 的地址
	// 这里暂时使用示例地址，需要根据实际情况修改
	apiServerURL := utils.GenAgentUrl(c.Config.Apiserver.Http.Scheme,
		c.Config.Apiserver.Http.Server,
		0,
		c.Config.Apiserver.Http.BaseUri,
		uriScriptResults)

	for _, task := range tasks {
		// 构建上报请求
		reportRequest := ReportScriptsExcute{
			TaskID:       task.TaskID,
			ScriptsID:    task.ScriptID,
			Output:       task.Output,
			Status:       task.Status,
			ErrorMessage: task.ErrorMsg,
		}

		// 发送请求到 API Server
		respBody, err := c.HTTPClient.Post(apiServerURL, reportRequest, nil)
		if err != nil {
			zap.L().Error("failed to report script execution result",
				zap.String("cron", "script_result_reporter"),
				zap.Uint("task_id", task.ID),
				zap.Uint("task_id_assigned", task.TaskID),
				zap.Error(err),
			)
			// 上报失败，不标记为已上报，下次继续尝试
			continue
		}

		// 解析响应
		var apiResp response.Response
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			zap.L().Error("failed to parse API Server response",
				zap.String("cron", "script_result_reporter"),
				zap.Uint("task_id", task.ID),
				zap.Uint("task_id_assigned", task.TaskID),
				zap.Error(err),
			)
			// 解析失败，不标记为已上报，下次继续尝试
			continue
		}

		// 检查响应码，Code=0 表示成功
		if apiResp.Code == 0 {
			// 只有 success 状态且上报成功时，才标记为已上报
			if task.Status == "success" {
				err := c.ScriptTaskRepo.MarkAsReported(task.ID)
				if err != nil {
					zap.L().Error("failed to mark task as reported",
						zap.String("cron", "script_result_reporter"),
						zap.Uint("task_id", task.ID),
						zap.Uint("task_id_assigned", task.TaskID),
						zap.Error(err),
					)
				} else {
					zap.L().Info("script execution result reported successfully",
						zap.String("cron", "script_result_reporter"),
						zap.Uint("task_id", task.ID),
						zap.Uint("task_id_assigned", task.TaskID),
						zap.String("status", task.Status),
					)
				}
			} else {
				// 非 success 状态，虽然上报成功但不标记，继续上报以保持状态同步
				zap.L().Info("script execution result reported (non-success status, will continue reporting)",
					zap.String("cron", "script_result_reporter"),
					zap.Uint("task_id", task.ID),
					zap.Uint("task_id_assigned", task.TaskID),
					zap.String("status", task.Status),
				)
			}
		} else {
			zap.L().Error("API Server returned error",
				zap.String("cron", "script_result_reporter"),
				zap.Uint("task_id", task.ID),
				zap.Uint("task_id_assigned", task.TaskID),
				zap.Int("code", apiResp.Code),
				zap.String("message", apiResp.Message),
			)
		}
	}
}
