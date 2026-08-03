package domain

import (
	"context"
	"time"
)

type Task struct {
	ID         uint
	ScriptID   uint
	TaskID     uint
	Name       string
	Content    string
	Status     string
	Output     string
	ErrorMsg   string
	ExecutedAt *time.Time
	Reported   bool
}

type Repository interface {
	Add(context.Context, *Task) error
	Get(context.Context, uint) (Task, error)
	GetRunningTask(context.Context) (Task, error)
	Update(context.Context, *Task) error
	GetUnreportedTasks(context.Context) ([]Task, error)
	MarkAsReported(context.Context, uint) error
}

type Executor interface {
	Execute(uint, string) (string, error)
}
