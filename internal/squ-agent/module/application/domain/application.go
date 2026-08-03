package domain

import "context"

const (
	StatusStarting = "starting"
	StatusRunning  = "running"
	StatusStopped  = "stopped"
	StatusUndeploy = "undeploy"
	StatusPaused   = "paused"
	StatusFailed   = "Failed"
)

type Application struct {
	ID          uint
	Name        string
	Description string
	Type        string
	OldStatus   string
	Status      string
	Content     string
	Version     string
	DeployID    uint64
	Env         []map[string]string
}

type Repository interface {
	List(context.Context) ([]Application, error)
	Get(context.Context, uint) (Application, error)
	GetByDeployID(context.Context, uint64) (Application, error)
	Delete(context.Context, uint) error
	Add(context.Context, *Application) error
	Update(context.Context, *Application) error
	Transaction(context.Context, func(Repository) error) error
}

type ConfigStore interface {
	Save(context.Context, string, string) error
}

type ComposeRuntime interface {
	DockerInstalled() bool
	ComposeAvailable() bool
	Prepare(uint64, string, []map[string]string) (string, string, error)
	ComposeFileExists(uint64) bool
	Up(string, string) error
	Start(string, string) error
	Stop(string, string) error
	Path(uint64) string
}
