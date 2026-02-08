package script

import (
	"encoding/json"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/internal/squ-apiserver/handler/script/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

// Execute 执行脚本
func (s *Script) Execute(request req.ExecuteScript) response.Response {
	// 1. 检查脚本是否存在
	script, err := s.Repository.Get(request.ScriptID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrScriptNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := s.ServerRepo.Get(request.ServerID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrServerNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 生成唯一的 TaskID
	taskID, err := utils.IDGenerate()
	if err != nil {
		zap.L().Error("生成 TaskID 失败",
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 4. 先在数据库中创建执行记录，状态为 running
	result := model.ScriptResult{
		TaskID:    taskID,
		ScriptID:  request.ScriptID,
		ServerID:  request.ServerID,
		ServerIP:  server.IpAddress,
		AgentPort: server.AgentPort,
		Status:    "running",
	}
	if err := s.Repository.AddScriptResult(&result); err != nil {
		zap.L().Error("创建脚本执行记录失败",
			zap.Uint64("task_id", taskID),
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 5. 构建发送给 Agent 的请求
	agentURL := utils.GenAgentUrl(s.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		s.Config.Agent.Http.BaseUrl,
		"script/execute")
	agentReq := s.modelToRequest(script)
	agentReq.TaskID = uint(taskID) // 将生成的 TaskID 传给 Agent

	respBody, err := s.HTTPClient.Post(agentURL, agentReq, nil)
	if err != nil {
		zap.L().Error("执行脚本请求发送失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		// 更新执行记录状态为 failed
		result.Status = "failed"
		result.ErrorMessage = "发送执行请求失败: " + err.Error()
		s.Repository.UpdateScriptResultByTaskID(taskID, &result)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	// 解析响应
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("解析 Agent 响应失败",
			zap.String("url", agentURL),
			zap.Error(err),
		)
		// 更新执行记录状态为 failed
		result.Status = "failed"
		result.ErrorMessage = "解析响应失败: " + err.Error()
		s.Repository.UpdateScriptResultByTaskID(taskID, &result)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("Agent 执行脚本失败",
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		// 更新执行记录状态为 failed
		result.Status = "failed"
		result.ErrorMessage = "Agent 返回错误: " + agentResp.Message
		s.Repository.UpdateScriptResultByTaskID(taskID, &result)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	return response.Success("脚本执行任务已提交")
}

// ReceiveScriptResult 接收脚本执行结果
func (s *Script) ReceiveScriptResult(request req.ScriptResultReport) response.Response {
	// 1. 验证脚本是否存在（可选）
	_, err := s.Repository.Get(request.ScriptID)
	if err != nil {
		zap.L().Warn("接收脚本执行结果：脚本不存在",
			zap.Uint("script_id", request.ScriptID),
		)
		// 继续执行，不返回错误
	}

	// 2. 根据 TaskID 直接更新执行记录
	result := model.ScriptResult{
		Output:       request.Output,
		Status:       request.Status,
		ErrorMessage: request.ErrorMessage,
	}
	if err := s.Repository.UpdateScriptResultByTaskID(uint64(request.TaskID), &result); err != nil {
		zap.L().Error("更新脚本执行结果失败",
			zap.Uint64("task_id", uint64(request.TaskID)),
			zap.Uint("script_id", request.ScriptID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	zap.L().Info("脚本执行结果已更新",
		zap.Uint64("task_id", uint64(request.TaskID)),
		zap.Uint("script_id", request.ScriptID),
		zap.String("status", request.Status),
	)

	return response.Success("success")
}

// GetResults 获取脚本执行结果
func (s *Script) GetResults(scriptID uint) response.Response {
	results, err := s.Repository.GetScriptResults(scriptID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	var resultRes []res.ScriptResult
	for _, r := range results {
		resultRes = append(resultRes, s.scriptResultToResponse(r))
	}

	return response.Success(resultRes)
}
