package application

import (
	"context"
	"errors"

	"squirrel-dev/internal/squ-agent/module/server/domain"
)

var ErrCollectorUnavailable = errors.New("host collector unavailable")

type Service struct {
	collector domain.HostInfoCollector
}

func NewService(collector domain.HostInfoCollector) *Service {
	return &Service{collector: collector}
}

func (s *Service) GetInfo(ctx context.Context) (*domain.HostInfo, error) {
	if s.collector == nil {
		return nil, ErrCollectorUnavailable
	}
	return s.collector.CollectHostInfo(ctx)
}
