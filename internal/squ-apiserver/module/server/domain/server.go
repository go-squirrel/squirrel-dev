package domain

import "context"

const (
	StatusOnline  = "online"
	StatusOffline = "offline"

	AuthTypePassword = "password"
	AuthTypeKey      = "privatekey"
)

type Server struct {
	ID            uint
	UUID          string
	Hostname      string
	IPAddress     string
	AgentPort     int
	SSHUsername   string
	SSHPassword   *string
	SSHPrivateKey *string
	SSHPassphrase *string
	SSHPort       int
	AuthType      string
	ServerAlias   *string
	Status        string
}

type Repository interface {
	List(context.Context) ([]Server, error)
	Get(context.Context, uint) (Server, error)
	Delete(context.Context, uint) error
	Add(context.Context, *Server) error
	Update(context.Context, *Server) error
	GetByUUID(context.Context, string) (Server, error)
}

type AgentInfoClient interface {
	GetInfo(context.Context, string, int) (string, map[string]any)
}

type SSHTester interface {
	Test(context.Context, Server) error
}
