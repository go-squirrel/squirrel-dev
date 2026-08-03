package application

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/script/domain"
)

type repositoryStub struct {
	script          domain.Script
	getErr          error
	addedResult     *domain.ScriptResult
	updatedResult   *domain.ScriptResult
	updatedTaskID   uint64
	updateResultErr error
}

func (r *repositoryStub) List(context.Context) ([]domain.Script, error) {
	return []domain.Script{r.script}, nil
}
func (r *repositoryStub) Get(context.Context, uint) (domain.Script, error) {
	return r.script, r.getErr
}
func (r *repositoryStub) Delete(context.Context, uint) error { return nil }
func (r *repositoryStub) Add(context.Context, *domain.Script) error {
	return nil
}
func (r *repositoryStub) Update(context.Context, *domain.Script) error {
	return nil
}
func (r *repositoryStub) AddResult(_ context.Context, result *domain.ScriptResult) error {
	copy := *result
	r.addedResult = &copy
	return nil
}
func (r *repositoryStub) ListResults(context.Context, uint) ([]domain.ScriptResult, error) {
	return nil, nil
}
func (r *repositoryStub) UpdateResultByTaskID(_ context.Context, taskID uint64, result *domain.ScriptResult) error {
	copy := *result
	r.updatedTaskID = taskID
	r.updatedResult = &copy
	return r.updateResultErr
}

type serverStub struct {
	server domain.Server
	err    error
}

func (s serverStub) Get(context.Context, uint) (domain.Server, error) { return s.server, s.err }

type agentStub struct{ err error }

func (a agentStub) Post(context.Context, domain.Server, string, any) error { return a.err }

type idStub struct {
	id  uint64
	err error
}

func (i idStub) Generate() (uint64, error) { return i.id, i.err }

func TestAddPreservesValidationAndTrimming(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository, serverStub{}, agentStub{}, idStub{})
	if _, err := service.Add(context.Background(), ScriptRequest{Name: "x", Content: " \n#!/bin/sh"}); !errors.Is(err, ErrInvalidContent) {
		t.Fatalf("leading whitespace must fail the legacy shebang check: %v", err)
	}
	if _, err := service.Add(context.Background(), ScriptRequest{Name: "x", Content: "#!/bin/sh\n "}); err != nil {
		t.Fatalf("valid script failed: %v", err)
	}
}

func TestExecuteAgentFailureMarksResultFailed(t *testing.T) {
	repository := &repositoryStub{script: domain.Script{ID: 2, Name: "test", Content: "#!/bin/sh"}}
	service := NewService(
		repository,
		serverStub{server: domain.Server{ID: 3, IPAddress: "127.0.0.1", AgentPort: 10750}},
		agentStub{err: errors.New("offline")},
		idStub{id: 42},
	)
	_, err := service.Execute(context.Background(), ExecuteRequest{ScriptID: 2, ServerID: 3})
	if !errors.Is(err, ErrExecuteFailed) {
		t.Fatalf("unexpected error: %v", err)
	}
	if repository.addedResult == nil || repository.addedResult.Status != "running" {
		t.Fatalf("running result not created: %#v", repository.addedResult)
	}
	if repository.updatedTaskID != 42 || repository.updatedResult.Status != "failed" ||
		repository.updatedResult.ErrorMessage != "agent execution failed: offline" {
		t.Fatalf("failed result not preserved: %#v", repository.updatedResult)
	}
}

func TestReceiveResultAllowsMissingScriptAndMissingTask(t *testing.T) {
	repository := &repositoryStub{getErr: gorm.ErrRecordNotFound}
	service := NewService(repository, serverStub{}, agentStub{}, idStub{})
	data, err := service.ReceiveResult(context.Background(), ResultReport{
		TaskID: 99, ScriptID: 404, Output: "done", Status: "success",
	})
	if err != nil || data != "success" {
		t.Fatalf("legacy callback should succeed: data=%q err=%v", data, err)
	}
	if repository.updatedTaskID != 99 || repository.updatedResult.Output != "done" {
		t.Fatalf("callback update not issued: %#v", repository.updatedResult)
	}
}
