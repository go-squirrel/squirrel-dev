package application

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/monitor/domain"
)

type serverStub struct{ err error }

func (s serverStub) Get(context.Context, uint) (domain.Server, error) {
	return domain.Server{ID: 7, IPAddress: "127.0.0.1", AgentPort: 10750}, s.err
}

type agentStub struct {
	path   string
	result domain.Result
	err    error
}

func (a *agentStub) Get(_ context.Context, _ domain.Server, path string) (domain.Result, error) {
	a.path = path
	return a.result, a.err
}

func TestRangePathAndAgentResponseArePreserved(t *testing.T) {
	agent := &agentStub{result: domain.Result{Message: "success", Data: map[string]any{"cpu": 1}}}
	service := NewService(serverStub{}, agent)
	result, err := service.BaseRange(7, "24h&sample=raw")
	if err != nil {
		t.Fatal(err)
	}
	if agent.path != "monitor/base?range=24h&sample=raw" {
		t.Fatalf("legacy raw range path changed: %q", agent.path)
	}
	if result.Message != "success" || result.Data == nil {
		t.Fatalf("agent response was not returned unchanged: %#v", result)
	}
}

func TestMonitorErrorMapping(t *testing.T) {
	service := NewService(serverStub{err: gorm.ErrRecordNotFound}, &agentStub{})
	if _, err := service.Stats(7); !errors.Is(err, ErrServerNotFound) {
		t.Fatalf("missing server error=%v", err)
	}
	service = NewService(serverStub{}, &agentStub{err: errors.New("offline")})
	if _, err := service.Stats(7); !errors.Is(err, ErrMonitorFailed) {
		t.Fatalf("agent failure error=%v", err)
	}
}
