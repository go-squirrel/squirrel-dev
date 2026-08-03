package application

import (
	"context"
	"errors"
	"time"

	"squirrel-dev/internal/squ-agent/module/script/domain"
)

var ErrAlreadyRunning = errors.New("script is already running")

type Request struct {
	ID      uint
	Name    string
	Content string
	TaskID  uint
}

type Service struct {
	repository domain.Repository
	executor   domain.Executor
}

func NewService(repository domain.Repository, executor domain.Executor) *Service {
	return &Service{repository: repository, executor: executor}
}

func (s *Service) Execute(ctx context.Context, request Request) error {
	if _, err := s.repository.GetRunningTask(ctx); err == nil {
		return ErrAlreadyRunning
	}
	task := domain.Task{
		ScriptID: request.ID, TaskID: request.TaskID, Name: request.Name,
		Content: request.Content, Status: "pending", Reported: false,
	}
	if err := s.repository.Add(ctx, &task); err != nil {
		return err
	}
	go s.executeAsync(task.ID, request.Content)
	return nil
}

func (s *Service) executeAsync(taskID uint, content string) {
	ctx := context.Background()
	task, err := s.repository.Get(ctx, taskID)
	if err != nil {
		return
	}
	task.Status = "running"
	executedAt := time.Now()
	task.ExecutedAt = &executedAt
	_ = s.repository.Update(ctx, &task)

	output, executeErr := s.executor.Execute(taskID, content)
	task.Output = output
	task.Status = "success"
	task.Reported = false
	task.ExecutedAt = &executedAt
	if executeErr != nil {
		task.Status = "failed"
		task.ErrorMsg = executeErr.Error()
	}
	_ = s.repository.Update(ctx, &task)
}
