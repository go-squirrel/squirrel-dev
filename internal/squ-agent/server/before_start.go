package server

import "go.uber.org/zap"

func (s *Server) beforeStart() {
	zap.L().Info("before start.")
}

func (s *Server) registerAgent() {
	zap.L().Info("run.")
}
