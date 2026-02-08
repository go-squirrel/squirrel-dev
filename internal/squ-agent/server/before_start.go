package server

import "go.uber.org/zap"

func (s *Server) beforeStart() {
	zap.L().Info("Starting agent server...")
}

func (s *Server) registerAgent() {
	zap.L().Info("Registering agent...")
}
