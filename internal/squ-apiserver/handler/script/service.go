package script

import (
	"encoding/json"
	"strings"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/internal/squ-apiserver/handler/script/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"

	scriptRepository "squirrel-dev/internal/squ-apiserver/repository/script"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
)

type Script struct {
	Config     *config.Config
	Repository scriptRepository.ScriptRepository
	ServerRepo serverRepository.Repository
	HTTPClient *httpclient.Client
}

func New(config *config.Config, scriptRepo scriptRepository.ScriptRepository, serverRepo serverRepository.Repository) *Script {
	hc := httpclient.NewClient(30 * time.Second)
	return &Script{
		Config:     config,
		Repository: scriptRepo,
		ServerRepo: serverRepo,
		HTTPClient: hc,
	}
}

func (s *Script) List() response.Response {
	var scripts []res.Script
	daoScripts, err := s.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoS := range daoScripts {
		scripts = append(scripts, res.Script{
			ID:      daoS.ID,
			Name:    daoS.Name,
			Content: daoS.Content,
		})
	}
	return response.Success(scripts)
}

func (s *Script) Get(id uint) response.Response {
	var scriptRes res.Script
	daoS, err := s.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	scriptRes = res.Script{
		ID:      daoS.ID,
		Name:    daoS.Name,
		Content: daoS.Content,
	}

	return response.Success(scriptRes)
}

func (s *Script) Delete(id uint) response.Response {
	err := s.Repository.Delete(id)
	if err != nil {
		zap.S().Error(err)
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Script) Add(request req.Script) response.Response {
	// 验证脚本名称和内容
	if request.Name == "" {
		zap.S().Error("script name is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	if request.Content == "" {
		zap.S().Error("script content is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 验证脚本内容是否以 shebang 开头
	if !strings.HasPrefix(request.Content, "#!") {
		zap.S().Error("script must start with shebang (#!)")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 清理脚本内容（去除首尾空白）
	request.Content = strings.TrimSpace(request.Content)

	modelReq := model.Script{
		Name:    request.Name,
		Content: request.Content,
	}

	err := s.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Script) Update(request req.Script) response.Response {
	// 验证脚本名称和内容
	if request.Name == "" {
		zap.S().Error("script name is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	if request.Content == "" {
		zap.S().Error("script content is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 验证脚本内容是否以 shebang 开头
	if !strings.HasPrefix(request.Content, "#!") {
		zap.S().Error("script must start with shebang (#!)")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 清理脚本内容（去除首尾空白）
	request.Content = strings.TrimSpace(request.Content)

	modelReq := model.Script{
		Name:    request.Name,
		Content: request.Content,
	}
	modelReq.ID = request.ID
	err := s.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

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
	agentReq := req.Script{
		ID:      script.ID,
		Name:    script.Name,
		Content: script.Content,
		TaskID:  uint(taskID), // 将生成的 TaskID 传给 Agent
	}

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
		resultRes = append(resultRes, res.ScriptResult{
			ID:           r.ID,
			TaskID:       r.TaskID,
			ScriptID:     r.ScriptID,
			ServerID:     r.ServerID,
			ServerIP:     r.ServerIP,
			AgentPort:    r.AgentPort,
			Output:       r.Output,
			Status:       r.Status,
			ErrorMessage: r.ErrorMessage,
			CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.Success(resultRes)
}
