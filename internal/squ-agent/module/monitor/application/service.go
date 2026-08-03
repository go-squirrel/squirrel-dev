package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"squirrel-dev/internal/squ-agent/module/monitor/domain"
)

var ErrInvalidTimeRange = errors.New("invalid time range")

type Service struct {
	cache      domain.Cache
	repository domain.Repository
	collector  domain.Collector
	now        func() time.Time
}

func NewService(cache domain.Cache, repository domain.Repository, collector domain.Collector) *Service {
	return &Service{
		cache:      cache,
		repository: repository,
		collector:  collector,
		now:        time.Now,
	}
}

func (s *Service) Stats(ctx context.Context) (domain.Stats, error) {
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, domain.StatsCacheKey); err == nil {
			if stats, ok := cached.(domain.Stats); ok {
				return stats, nil
			}
		}
	}
	if s.collector == nil {
		return domain.Stats{}, fmt.Errorf("collector is nil")
	}
	stats, err := s.collector.Stats(ctx)
	if err != nil {
		return domain.Stats{}, err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, domain.StatsCacheKey, stats, domain.StatsCacheTTL)
	}
	return stats, nil
}

func (s *Service) AllDiskIO(ctx context.Context) (domain.AllDiskIOStats, error) {
	if s.collector == nil {
		return domain.AllDiskIOStats{}, fmt.Errorf("collector is nil")
	}
	return s.collector.AllDiskIO(ctx)
}

func (s *Service) DiskIO(ctx context.Context, device string) (domain.DiskIOStats, error) {
	if s.collector == nil {
		return domain.DiskIOStats{}, fmt.Errorf("collector is nil")
	}
	return s.collector.DiskIO(ctx, device)
}

func (s *Service) AllNetIO(ctx context.Context) (domain.AllNetIOStats, error) {
	if s.collector == nil {
		return domain.AllNetIOStats{}, fmt.Errorf("collector is nil")
	}
	return s.collector.AllNetIO(ctx)
}

func (s *Service) NetIO(ctx context.Context, interfaceName string) (domain.NetIOStats, error) {
	if s.collector == nil {
		return domain.NetIOStats{}, fmt.Errorf("collector is nil")
	}
	return s.collector.NetIO(ctx, interfaceName)
}

func (s *Service) BaseByRange(ctx context.Context, value string) ([]domain.BaseMonitor, error) {
	since, err := s.parseTimeRange(value)
	if err != nil {
		return nil, err
	}
	return s.repository.BaseByTimeRange(ctx, since)
}

func (s *Service) DiskIOByRange(ctx context.Context, value string) ([]domain.DiskIOMonitor, error) {
	since, err := s.parseTimeRange(value)
	if err != nil {
		return nil, err
	}
	return s.repository.DiskIOByTimeRange(ctx, since)
}

func (s *Service) DiskUsageByRange(ctx context.Context, value string) ([]domain.DiskUsageMonitor, error) {
	since, err := s.parseTimeRange(value)
	if err != nil {
		return nil, err
	}
	return s.repository.DiskUsageByTimeRange(ctx, since)
}

func (s *Service) NetworkByRange(ctx context.Context, value string) ([]domain.NetworkMonitor, error) {
	since, err := s.parseTimeRange(value)
	if err != nil {
		return nil, err
	}
	return s.repository.NetworkByTimeRange(ctx, since)
}

func (s *Service) parseTimeRange(value string) (time.Time, error) {
	switch value {
	case "1h":
		return s.now().Add(-time.Hour), nil
	case "6h":
		return s.now().Add(-6 * time.Hour), nil
	case "24h":
		return s.now().Add(-24 * time.Hour), nil
	case "7d":
		return s.now().Add(-7 * 24 * time.Hour), nil
	default:
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidTimeRange, value)
	}
}
