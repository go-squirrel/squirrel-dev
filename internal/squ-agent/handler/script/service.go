package script

import (
	"fmt"
	"os"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	scriptReq "squirrel-dev/internal/squ-agent/handler/script/req"
	scriptRes "squirrel-dev/internal/squ-agent/handler/script/res"
	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/execute"

	scriptTaskRepo "squirrel-dev/internal/squ-agent/repository/script_task"

	"go.uber.org/zap"
)

type Script struct {
	Config   *config.Config
	TaskRepo scriptTaskRepo.Repository
}

func New(config *config.Config, taskRepo scriptTaskRepo.Repository) *Script {
	return &Script{
		Config:   config,
		TaskRepo: taskRepo,
	}
}

// Execute 执行脚本
func (s *Script) Execute(request scriptReq.Script) response.Response {
	// 1. 检查是否有正在运行的脚本
	runningTask, err := s.TaskRepo.GetRunningTask()
	if err == nil {
		zap.L().Warn("Script already running",
			zap.Uint("task_id", runningTask.ID),
			zap.String("name", runningTask.Name),
		)
		return response.Error(scriptRes.ErrScriptAlreadyRunning)
	}

	// 2. 创建任务记录（TaskID 由 APIServer 分配）
	task := model.ScriptExecutionTask{
		ScriptID: request.ID,
		TaskID:   request.TaskID, // 保存 APIServer 分配的 TaskID
		Name:     request.Name,
		Content:  request.Content,
		Status:   "pending",
		Reported: false,
	}

	err = s.TaskRepo.Add(&task)
	if err != nil {
		zap.L().Error("Failed to create script execution task",
			zap.Uint("task_id", request.TaskID),
			zap.String("name", request.Name),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 异步执行脚本
	go s.executeScriptAsync(task.ID, request.Content)

	return response.Success("Script execution task created")
}

// executeScriptAsync 异步执行脚本
func (s *Script) executeScriptAsync(taskID uint, content string) {
	// 更新任务状态为 running
	task, err := s.TaskRepo.Get(taskID)
	if err != nil {
		zap.L().Error("Failed to get task",
			zap.Uint("task_id", taskID),
			zap.Error(err),
		)
		return
	}

	task.Status = "running"
	executedAt := time.Now()
	task.ExecutedAt = &executedAt
	err = s.TaskRepo.Update(&task)
	if err != nil {
		zap.L().Error("Failed to update task status",
			zap.Uint("task_id", taskID),
			zap.Error(err),
		)
	}

	// 创建临时脚本文件
	tmpFile := fmt.Sprintf("/tmp/script_%d.sh", taskID)
	err = os.WriteFile(tmpFile, []byte(content), 0755)
	if err != nil {
		zap.L().Error("Failed to create temporary script file",
			zap.Uint("task_id", taskID),
			zap.String("file", tmpFile),
			zap.Error(err),
		)
		s.updateTaskFailed(taskID, err.Error())
		return
	}
	defer os.Remove(tmpFile)

	// 执行脚本
	output, err := execute.Command("bash", tmpFile)
	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = err.Error()
		zap.L().Error("Failed to execute script",
			zap.Uint("task_id", taskID),
			zap.String("error", err.Error()),
		)
	}
	// 更新任务结果
	task.Output = string(output)
	task.Status = "success"
	task.Reported = false
	task.ExecutedAt = &executedAt

	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = err.Error()
		zap.L().Error("Failed to execute script",
			zap.Uint("task_id", taskID),
			zap.String("error", err.Error()),
		)
	} else {
		zap.L().Info("Script executed successfully",
			zap.Uint("task_id", taskID),
			zap.Int("output_length", len(output)),
		)
	}

	err = s.TaskRepo.Update(&task)
	if err != nil {
		zap.L().Error("Failed to update task result",
			zap.Uint("task_id", taskID),
			zap.Error(err),
		)
	}
}

// updateTaskFailed 更新任务为失败状态
func (s *Script) updateTaskFailed(taskID uint, errorMsg string) {
	task, err := s.TaskRepo.Get(taskID)
	if err != nil {
		return
	}

	task.Status = "failed"
	task.ErrorMsg = errorMsg
	task.Reported = false

	err = s.TaskRepo.Update(&task)
	if err != nil {
		zap.L().Error("Failed to update task failed status",
			zap.Uint("task_id", taskID),
			zap.Error(err),
		)
	}
}
