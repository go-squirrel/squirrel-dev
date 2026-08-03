package domain

import (
	"context"
)

type Server struct {
	ID        uint
	IPAddress string
	AgentPort int
}

type Result struct {
	Message string
	Data    any
}

type ServerReader interface {
	Get(context.Context, uint) (Server, error)
}

type AgentClient interface {
	Get(context.Context, Server, string) (Result, error)
}
