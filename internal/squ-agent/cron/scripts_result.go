package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// startScriptResultReporter 启动脚本执行结果上报定时任务
// 每 5 秒上报一次未上报的执行结果
func (c *Cron) startScriptResultReporter() error {
	_, err := c.Cron.AddFunc("*/5 * * * * *", func() {
		c.reportScriptResults()
	})
	if err != nil {
		zap.L().Error("启动脚本结果上报定时任务失败", zap.Error(err))
	}
	return err
}

// reportScriptResults 上报脚本执行结果到 API Server
func (c *Cron) reportScriptResults() {
	// 获取未上报的任务
	tasks, err := c.ScriptTaskRepo.GetUnreportedTasks()
	if err != nil {
		zap.L().Error("获取未上报的脚本任务失败", zap.Error(err))
		return
	}

	if len(tasks) == 0 {
		return
	}

	zap.L().Info("开始上报脚本执行结果", zap.Int("count", len(tasks)))
	// TODO: 从配置中获取 API Server 的地址
	// 这里暂时使用示例地址，需要根据实际情况修改
	apiServerURL := fmt.Sprintf("%s%s", c.Config.Apiserver.Url, uriScriptResults)

	for _, task := range tasks {
		// 构建上报请求
		reportRequest := map[string]any{
			"script_id":     task.ScriptID,
			"server_id":     1,           // TODO: 需要从配置中获取本服务器的 ID
			"server_ip":     "127.0.0.1", // TODO: 需要从配置中获取
			"agent_port":    9527,        // TODO: 需要从配置中获取
			"output":        task.Output,
			"status":        task.Status,
			"error_message": task.ErrorMsg,
		}

		// 发送请求到 API Server
		jsonData, _ := json.Marshal(reportRequest)
		req, err := http.NewRequest("POST", apiServerURL, bytes.NewBuffer(jsonData))
		if err != nil {
			zap.L().Error("创建上报请求失败",
				zap.Uint("task_id", task.ID),
				zap.Error(err),
			)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			zap.L().Error("上报脚本执行结果失败",
				zap.Uint("task_id", task.ID),
				zap.Error(err),
			)
			continue
		}
		resp.Body.Close()

		// 如果上报成功，标记任务为已上报
		if resp.StatusCode == http.StatusOK {
			err := c.ScriptTaskRepo.MarkAsReported(task.ID)
			if err != nil {
				zap.L().Error("标记任务为已上报失败",
					zap.Uint("task_id", task.ID),
					zap.Error(err),
				)
			} else {
				zap.L().Info("脚本执行结果上报成功",
					zap.Uint("task_id", task.ID),
					zap.String("status", task.Status),
				)
			}
		} else {
			zap.L().Error("上报脚本执行结果失败",
				zap.Uint("task_id", task.ID),
				zap.Int("status_code", resp.StatusCode),
			)
		}
	}
}
