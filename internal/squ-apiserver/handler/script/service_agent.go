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
	// 1. Check if script exists
	script, err := s.Repository.Get(request.ScriptID)
	if err != nil {
		zap.L().Error("failed to get script for execution",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	// 2. Check if server exists
	server, err := s.ServerRepo.Get(request.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for script execution",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrServerNotFound)
	}

	// 3. Generate unique TaskID
	taskID, err := utils.IDGenerate()
	if err != nil {
		zap.L().Error("failed to generate task ID",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	// 4. Create execution record in database with status running
	result := model.ScriptResult{
		TaskID:    taskID,
		ScriptID:  request.ScriptID,
		ServerID:  request.ServerID,
		ServerIP:  server.IpAddress,
		AgentPort: server.AgentPort,
		Status:    "running",
	}
	if err := s.Repository.AddScriptResult(&result); err != nil {
		zap.L().Error("failed to create script execution record",
			zap.Uint64("task_id", taskID),
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	// 5. Build request for Agent
	agentURL := utils.GenAgentUrl(s.Config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		s.Config.Agent.Http.BaseUrl,
		"script/execute")
	agentReq := s.modelToRequest(script)
	agentReq.TaskID = uint(taskID) // Pass generated TaskID to Agent

	respBody, err := s.HTTPClient.Post(agentURL, agentReq, nil)
	if err != nil {
		zap.L().Error("failed to send script execution request",
			zap.Uint64("task_id", taskID),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		// Update execution record status to failed
		result.Status = "failed"
		result.ErrorMessage = "failed to send execution request: " + err.Error()
		s.Repository.UpdateScriptResultByTaskID(taskID, &result)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	// Parse response
	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("failed to parse agent response",
			zap.Uint64("task_id", taskID),
			zap.String("url", agentURL),
			zap.Error(err),
		)
		// Update execution record status to failed
		result.Status = "failed"
		result.ErrorMessage = "failed to parse response: " + err.Error()
		s.Repository.UpdateScriptResultByTaskID(taskID, &result)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	if agentResp.Code != 0 {
		zap.L().Error("agent failed to execute script",
			zap.Uint64("task_id", taskID),
			zap.String("url", agentURL),
			zap.Int("code", agentResp.Code),
			zap.String("message", agentResp.Message),
		)
		// Update execution record status to failed
		result.Status = "failed"
		result.ErrorMessage = "agent returned error: " + agentResp.Message
		s.Repository.UpdateScriptResultByTaskID(taskID, &result)
		return response.Error(res.ErrScriptExecutionFailed)
	}

	return response.Success("script execution task submitted")
}

// ReceiveScriptResult 接收脚本执行结果
func (s *Script) ReceiveScriptResult(request req.ScriptResultReport) response.Response {
	// 1. Validate script exists (optional)
	_, err := s.Repository.Get(request.ScriptID)
	if err != nil {
		zap.L().Warn("script not found when receiving execution result",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint64("task_id", uint64(request.TaskID)),
			zap.Error(err),
		)
		// Continue execution, do not return error
	}

	// 2. Update execution record by TaskID
	result := model.ScriptResult{
		Output:       request.Output,
		Status:       request.Status,
		ErrorMessage: request.ErrorMessage,
	}
	if err := s.Repository.UpdateScriptResultByTaskID(uint64(request.TaskID), &result); err != nil {
		zap.L().Error("failed to update script execution result",
			zap.Uint64("task_id", uint64(request.TaskID)),
			zap.Uint("script_id", request.ScriptID),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	zap.L().Info("script execution result updated",
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
		zap.L().Error("failed to get script execution results",
			zap.Uint("script_id", scriptID),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	var resultRes []res.ScriptResult
	for _, r := range results {
		resultRes = append(resultRes, s.scriptResultToResponse(r))
	}

	return response.Success(resultRes)
}
