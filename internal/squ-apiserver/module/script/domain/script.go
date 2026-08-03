package domain

import (
	"context"
	"time"
)

type Script struct {
	ID      uint
	Name    string
	Content string
}

type ScriptResult struct {
	ID           uint
	CreatedAt    time.Time
	TaskID       uint64
	ScriptID     uint
	ServerID     uint
	ServerIP     string
	AgentPort    int
	Output       string
	Status       string
	ErrorMessage string
}

type Server struct {
	ID        uint
	IPAddress string
	AgentPort int
}

type Repository interface {
	List(context.Context) ([]Script, error)
	Get(context.Context, uint) (Script, error)
	Delete(context.Context, uint) error
	Add(context.Context, *Script) error
	Update(context.Context, *Script) error
	AddResult(context.Context, *ScriptResult) error
	ListResults(context.Context, uint) ([]ScriptResult, error)
	UpdateResultByTaskID(context.Context, uint64, *ScriptResult) error
}

type ServerReader interface {
	Get(context.Context, uint) (Server, error)
}

type AgentClient interface {
	Post(context.Context, Server, string, any) error
}

type IDGenerator interface {
	Generate() (uint64, error)
}
