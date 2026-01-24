package script

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	scriptReq "squirrel-dev/internal/squ-agent/handler/script/req"
	scriptRes "squirrel-dev/internal/squ-agent/handler/script/res"
	"squirrel-dev/internal/squ-agent/model"

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
		zap.L().Warn("已有脚本正在执行，等待中...",
			zap.Uint("task_id", runningTask.ID),
			zap.String("name", runningTask.Name),
		)
		return response.Error(scriptRes.ErrScriptAlreadyRunning)
	}

	// 2. 创建任务记录
	task := model.ScriptExecutionTask{
		ScriptID: request.ID,
		Name:     request.Name,
		Content:  request.Content,
		Status:   "pending",
		Reported: false,
	}

	err = s.TaskRepo.Add(&task)
	if err != nil {
		zap.L().Error("创建脚本执行任务失败",
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 异步执行脚本
	go s.executeScriptAsync(task.ID, request.Content)

	return response.Success("脚本执行任务已创建")
}

// executeScriptAsync 异步执行脚本
func (s *Script) executeScriptAsync(taskID uint, content string) {
	// 更新任务状态为 running
	task, err := s.TaskRepo.Get(taskID)
	if err != nil {
		zap.L().Error("获取任务失败",
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
		zap.L().Error("更新任务状态失败",
			zap.Uint("task_id", taskID),
			zap.Error(err),
		)
	}

	// 创建临时脚本文件
	tmpFile := fmt.Sprintf("/tmp/script_%d.sh", taskID)
	err = os.WriteFile(tmpFile, []byte(content), 0755)
	if err != nil {
		zap.L().Error("创建临时脚本文件失败",
			zap.Uint("task_id", taskID),
			zap.String("file", tmpFile),
			zap.Error(err),
		)
		s.updateTaskFailed(taskID, err.Error())
		return
	}
	defer os.Remove(tmpFile)

	// 执行脚本
	cmd := exec.Command("bash", tmpFile)
	output, err := cmd.CombinedOutput()

	// 更新任务结果
	task, _ = s.TaskRepo.Get(taskID)
	task.Output = string(output)
	task.Status = "success"
	task.Reported = false
	task.ExecutedAt = &executedAt

	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = err.Error()
		zap.L().Error("脚本执行失败",
			zap.Uint("task_id", taskID),
			zap.String("error", err.Error()),
		)
	} else {
		zap.L().Info("脚本执行成功",
			zap.Uint("task_id", taskID),
			zap.Int("output_length", len(output)),
		)
	}

	err = s.TaskRepo.Update(&task)
	if err != nil {
		zap.L().Error("更新任务结果失败",
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
		zap.L().Error("更新任务失败状态错误",
			zap.Uint("task_id", taskID),
			zap.Error(err),
		)
	}
}
