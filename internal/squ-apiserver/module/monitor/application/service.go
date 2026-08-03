package application

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/monitor/domain"
)

type Service struct {
	servers domain.ServerReader
	agent   domain.AgentClient
}

func NewService(servers domain.ServerReader, agent domain.AgentClient) *Service {
	return &Service{servers: servers, agent: agent}
}

func (s *Service) Stats(serverID uint) (domain.Result, error) {
	return s.callAgent(serverID, "monitor/stats")
}

func (s *Service) DiskIO(serverID uint, device string) (domain.Result, error) {
	return s.callAgent(serverID, fmt.Sprintf("monitor/stats/io/%s", device))
}

func (s *Service) AllDiskIO(serverID uint) (domain.Result, error) {
	return s.callAgent(serverID, "monitor/stats/io/all")
}

func (s *Service) NetIO(serverID uint, interfaceName string) (domain.Result, error) {
	return s.callAgent(serverID, fmt.Sprintf("monitor/stats/net/%s", interfaceName))
}

func (s *Service) AllNetIO(serverID uint) (domain.Result, error) {
	return s.callAgent(serverID, "monitor/stats/net/all")
}

func (s *Service) BaseRange(serverID uint, timeRange string) (domain.Result, error) {
	return s.callAgent(serverID, fmt.Sprintf("monitor/base?range=%s", timeRange))
}

func (s *Service) DiskRange(serverID uint, timeRange string) (domain.Result, error) {
	return s.callAgent(serverID, fmt.Sprintf("monitor/disk?range=%s", timeRange))
}

func (s *Service) DiskUsageRange(serverID uint, timeRange string) (domain.Result, error) {
	return s.callAgent(serverID, fmt.Sprintf("monitor/disk-usage?range=%s", timeRange))
}

func (s *Service) NetworkRange(serverID uint, timeRange string) (domain.Result, error) {
	return s.callAgent(serverID, fmt.Sprintf("monitor/net?range=%s", timeRange))
}

func (s *Service) callAgent(serverID uint, path string) (domain.Result, error) {
	// Legacy monitor calls were detached from the incoming HTTP request.
	server, err := s.servers.Get(context.Background(), serverID)
	if err != nil {
		zap.L().Error("failed to get server for monitoring",
			zap.Uint("server_id", serverID),
			zap.String("agent_path", path),
			zap.Error(err),
		)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Result{}, ErrServerNotFound
		}
		return domain.Result{}, ErrMonitorFailed
	}
	result, err := s.agent.Get(context.Background(), server, path)
	if err != nil {
		zap.L().Error("failed to get monitor data from agent",
			zap.Uint("server_id", serverID),
			zap.String("agent_path", path),
			zap.Error(err),
		)
		return domain.Result{}, ErrMonitorFailed
	}
	return result, nil
}
