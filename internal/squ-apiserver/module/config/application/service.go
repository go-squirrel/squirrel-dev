package application

import (
	"context"

	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/config/domain"
)

type Service struct{ repository domain.Repository }

func NewService(repository domain.Repository) *Service { return &Service{repository: repository} }

func (s *Service) List(ctx context.Context) ([]domain.Config, error) {
	configs, err := s.repository.List(ctx)
	if err != nil {
		zap.L().Error("failed to list configs", zap.Error(err))
		return nil, err
	}
	return configs, nil
}

func (s *Service) Get(ctx context.Context, id uint) (domain.Config, error) {
	config, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get config", zap.Uint("config_id", id), zap.Error(err))
		return domain.Config{}, err
	}
	return config, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		zap.L().Error("failed to delete config", zap.Uint("config_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Add(ctx context.Context, key, value string) error {
	if err := s.repository.Add(ctx, &domain.Config{Key: key, Value: value}); err != nil {
		zap.L().Error("failed to add config", zap.String("config_key", key), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, id uint, key, value string) error {
	if err := s.repository.Update(ctx, &domain.Config{ID: id, Key: key, Value: value}); err != nil {
		zap.L().Error("failed to update config",
			zap.Uint("config_id", id),
			zap.String("config_key", key),
			zap.Error(err),
		)
		return err
	}
	return nil
}
