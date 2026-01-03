package server

import (
	"squirrel-dev/internal/squctl/config"
)

type Server struct {
	Config *config.Config
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	return
}
