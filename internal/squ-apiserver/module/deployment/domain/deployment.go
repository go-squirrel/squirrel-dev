package domain

import (
	"context"
	"time"
)

type Deployment struct {
	ID            uint
	CreatedAt     time.Time
	ServerID      uint
	ApplicationID uint
	Status        string
	DeployID      uint64
	Content       string
	Env           []map[string]string
}

type Application struct {
	ID          uint
	Name        string
	Description string
	Type        string
	Content     string
	Version     string
}

type Server struct {
	ID        uint
	IPAddress string
	AgentPort int
}

type Repository interface {
	List(context.Context, uint) ([]Deployment, error)
	Get(context.Context, uint) (Deployment, error)
	GetByDeployID(context.Context, uint64) (Deployment, error)
	Delete(context.Context, uint) error
	Add(context.Context, *Deployment) error
	Update(context.Context, *Deployment) error
	UpdateStatus(context.Context, uint64, string) error
}

type ApplicationReader interface {
	Get(context.Context, uint) (Application, error)
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
