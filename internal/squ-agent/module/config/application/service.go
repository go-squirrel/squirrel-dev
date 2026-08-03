package application

import (
	"context"

	"squirrel-dev/internal/squ-agent/module/config/domain"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]domain.Config, error) {
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id uint) (domain.Config, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Save(ctx context.Context, key, value string) error {
	return s.repository.CreateOrUpdate(ctx, &domain.Config{
		Key:   key,
		Value: value,
	})
}

func (s *Service) GetByKey(ctx context.Context, key string) (string, error) {
	return s.repository.GetByKey(ctx, key)
}
