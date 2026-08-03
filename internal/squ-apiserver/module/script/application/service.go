package application

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/script/domain"
)

type ScriptRequest struct {
	ID      uint
	Name    string
	Content string
}

type ExecuteRequest struct {
	ScriptID uint
	ServerID uint
}

type ResultReport struct {
	TaskID       uint
	ScriptID     uint
	Output       string
	Status       string
	ErrorMessage string
}

type AgentScriptRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TaskID  uint   `json:"task_id"`
}

type Service struct {
	repository domain.Repository
	servers    domain.ServerReader
	agent      domain.AgentClient
	ids        domain.IDGenerator
}

func NewService(
	repository domain.Repository,
	servers domain.ServerReader,
	agent domain.AgentClient,
	ids domain.IDGenerator,
) *Service {
	return &Service{repository: repository, servers: servers, agent: agent, ids: ids}
}

func (s *Service) List(ctx context.Context) ([]domain.Script, error) {
	values, err := s.repository.List(ctx)
	if err != nil {
		zap.L().Error("failed to list scripts", zap.Error(err))
		return nil, repositoryError(err)
	}
	return values, nil
}

func (s *Service) Get(ctx context.Context, id uint) (domain.Script, error) {
	value, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get script", zap.Uint("script_id", id), zap.Error(err))
		return domain.Script{}, repositoryError(err)
	}
	return value, nil
}

func (s *Service) Delete(ctx context.Context, id uint) (string, error) {
	if err := s.repository.Delete(ctx, id); err != nil {
		zap.L().Error("failed to delete script", zap.Uint("script_id", id), zap.Error(err))
		return "", repositoryError(err)
	}
	return "success", nil
}

func (s *Service) Add(ctx context.Context, request ScriptRequest) (string, error) {
	if request.Name == "" {
		zap.L().Warn("script name is empty", zap.String("operation", "add"))
		return "", ErrInvalidContent
	}
	if request.Content == "" {
		zap.L().Warn("script content is empty", zap.String("script_name", request.Name))
		return "", ErrInvalidContent
	}
	if !strings.HasPrefix(request.Content, "#!") {
		zap.L().Warn("script must start with a shebang", zap.String("script_name", request.Name))
		return "", ErrInvalidContent
	}
	value := domain.Script{Name: request.Name, Content: strings.TrimSpace(request.Content)}
	if err := s.repository.Add(ctx, &value); err != nil {
		zap.L().Error("failed to add script", zap.String("script_name", request.Name), zap.Error(err))
		return "", repositoryError(err)
	}
	return "success", nil
}

func (s *Service) Update(ctx context.Context, request ScriptRequest) (string, error) {
	if request.Name == "" {
		zap.L().Warn("script name is empty", zap.Uint("script_id", request.ID))
		return "", ErrInvalidContent
	}
	if request.Content == "" {
		zap.L().Warn("script content is empty",
			zap.Uint("script_id", request.ID),
			zap.String("script_name", request.Name),
		)
		return "", ErrInvalidContent
	}
	if !strings.HasPrefix(request.Content, "#!") {
		zap.L().Warn("script must start with a shebang",
			zap.Uint("script_id", request.ID),
			zap.String("script_name", request.Name),
		)
		return "", ErrInvalidContent
	}
	value := domain.Script{
		ID: request.ID, Name: request.Name, Content: strings.TrimSpace(request.Content),
	}
	if err := s.repository.Update(ctx, &value); err != nil {
		zap.L().Error("failed to update script",
			zap.Uint("script_id", request.ID),
			zap.String("script_name", request.Name),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	return "success", nil
}

func (s *Service) Execute(ctx context.Context, request ExecuteRequest) (string, error) {
	script, err := s.repository.Get(ctx, request.ScriptID)
	if err != nil {
		zap.L().Error("failed to get script for execution",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	server, err := s.servers.Get(ctx, request.ServerID)
	if err != nil {
		zap.L().Error("failed to get server for script execution",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return "", ErrServerNotFound
	}
	taskID, err := s.ids.Generate()
	if err != nil {
		zap.L().Error("failed to generate script task ID",
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return "", ErrExecuteFailed
	}
	result := domain.ScriptResult{
		TaskID: taskID, ScriptID: request.ScriptID, ServerID: request.ServerID,
		ServerIP: server.IPAddress, AgentPort: server.AgentPort, Status: "running",
	}
	if err := s.repository.AddResult(ctx, &result); err != nil {
		zap.L().Error("failed to create script execution record",
			zap.Uint64("task_id", taskID),
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	agentRequest := AgentScriptRequest{
		ID: script.ID, Name: script.Name, Content: script.Content, TaskID: uint(taskID),
	}
	if err := s.agent.Post(ctx, server, "script/execute", agentRequest); err != nil {
		zap.L().Error("failed to submit script execution to agent",
			zap.Uint64("task_id", taskID),
			zap.Uint("script_id", request.ScriptID),
			zap.Uint("server_id", request.ServerID),
			zap.Error(err),
		)
		result.Status = "failed"
		result.ErrorMessage = "agent execution failed: " + err.Error()
		if updateErr := s.repository.UpdateResultByTaskID(ctx, taskID, &result); updateErr != nil {
			zap.L().Error("failed to mark script execution as failed",
				zap.Uint64("task_id", taskID),
				zap.Uint("script_id", request.ScriptID),
				zap.Uint("server_id", request.ServerID),
				zap.Error(updateErr),
			)
		}
		return "", ErrExecuteFailed
	}
	return "script execution task submitted", nil
}

func (s *Service) ReceiveResult(ctx context.Context, request ResultReport) (string, error) {
	// The old service deliberately accepted reports for a missing script.
	if _, err := s.repository.Get(ctx, request.ScriptID); err != nil {
		zap.L().Warn("failed to get script when receiving execution result",
			zap.Uint64("task_id", uint64(request.TaskID)),
			zap.Uint("script_id", request.ScriptID),
			zap.Error(err),
		)
	}
	result := domain.ScriptResult{
		Output: request.Output, Status: request.Status, ErrorMessage: request.ErrorMessage,
	}
	if err := s.repository.UpdateResultByTaskID(ctx, uint64(request.TaskID), &result); err != nil {
		zap.L().Error("failed to update script execution result",
			zap.Uint64("task_id", uint64(request.TaskID)),
			zap.Uint("script_id", request.ScriptID),
			zap.String("status", request.Status),
			zap.Error(err),
		)
		return "", repositoryError(err)
	}
	zap.L().Info("script execution result updated",
		zap.Uint64("task_id", uint64(request.TaskID)),
		zap.Uint("script_id", request.ScriptID),
		zap.String("status", request.Status),
	)
	return "success", nil
}

func (s *Service) ListResults(ctx context.Context, scriptID uint) ([]domain.ScriptResult, error) {
	values, err := s.repository.ListResults(ctx, scriptID)
	if err != nil {
		zap.L().Error("failed to list script execution results", zap.Uint("script_id", scriptID), zap.Error(err))
		return nil, repositoryError(err)
	}
	return values, nil
}
