package application

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-agent/module/script/domain"
)

type memoryRepository struct {
	mu     sync.Mutex
	nextID uint
	tasks  map[uint]domain.Task
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{nextID: 1, tasks: make(map[uint]domain.Task)}
}

func (r *memoryRepository) Add(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task.ID = r.nextID
	r.nextID++
	r.tasks[task.ID] = *task
	return nil
}
func (r *memoryRepository) Get(_ context.Context, id uint) (domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[id]
	if !ok {
		return domain.Task{}, gorm.ErrRecordNotFound
	}
	return task, nil
}
func (r *memoryRepository) GetRunningTask(context.Context) (domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, task := range r.tasks {
		if task.Status == "running" {
			return task, nil
		}
	}
	return domain.Task{}, gorm.ErrRecordNotFound
}
func (r *memoryRepository) Update(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = *task
	return nil
}
func (r *memoryRepository) GetUnreportedTasks(context.Context) ([]domain.Task, error) {
	return nil, nil
}
func (r *memoryRepository) MarkAsReported(context.Context, uint) error { return nil }

type fakeExecutor struct {
	output string
	err    error
	done   chan struct{}
}

func (e fakeExecutor) Execute(uint, string) (string, error) {
	close(e.done)
	return e.output, e.err
}

func TestExecuteAsyncKeepsLegacyTaskLifecycle(t *testing.T) {
	repository := newMemoryRepository()
	done := make(chan struct{})
	service := NewService(repository, fakeExecutor{output: "done", done: done})

	err := service.Execute(context.Background(), Request{ID: 7, TaskID: 9, Name: "demo", Content: "echo done"})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("script executor did not run")
	}

	deadline := time.Now().Add(time.Second)
	for {
		task, err := repository.Get(context.Background(), 1)
		if err != nil {
			t.Fatal(err)
		}
		if task.Status == "success" {
			if task.ScriptID != 7 || task.TaskID != 9 || task.Output != "done" || task.Reported || task.ExecutedAt == nil {
				t.Fatalf("unexpected completed task: %#v", task)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("task did not complete: %#v", task)
		}
	}
}

func TestExecuteRejectsRunningTask(t *testing.T) {
	repository := newMemoryRepository()
	repository.tasks[1] = domain.Task{ID: 1, Status: "running"}
	service := NewService(repository, fakeExecutor{done: make(chan struct{})})
	if err := service.Execute(context.Background(), Request{}); !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("err = %v", err)
	}
}
