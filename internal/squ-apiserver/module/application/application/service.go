package application

import (
	"context"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"squirrel-dev/internal/squ-apiserver/module/application/domain"
)

type Service struct{ repository domain.Repository }

func NewService(repository domain.Repository) *Service { return &Service{repository: repository} }

func (s *Service) List(ctx context.Context) ([]domain.Application, error) {
	values, err := s.repository.List(ctx)
	if err != nil {
		zap.L().Error("failed to list applications", zap.Error(err))
		return nil, err
	}
	return values, nil
}

func (s *Service) Get(ctx context.Context, id uint) (domain.Application, error) {
	value, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get application", zap.Uint("application_id", id), zap.Error(err))
		return domain.Application{}, err
	}
	return value, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		zap.L().Error("failed to delete application", zap.Uint("application_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Add(ctx context.Context, value domain.Application) error {
	if err := validateYAML(value.Content); err != nil {
		zap.L().Warn("invalid application YAML", zap.String("application_name", value.Name), zap.Error(err))
		return ErrInvalidYAML
	}
	if err := s.repository.Add(ctx, &value); err != nil {
		zap.L().Error("failed to add application", zap.String("application_name", value.Name), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, value domain.Application) error {
	if err := validateYAML(value.Content); err != nil {
		zap.L().Warn("invalid application YAML",
			zap.Uint("application_id", value.ID),
			zap.String("application_name", value.Name),
			zap.Error(err),
		)
		return ErrInvalidYAML
	}
	if err := s.repository.Update(ctx, &value); err != nil {
		zap.L().Error("failed to update application",
			zap.Uint("application_id", value.ID),
			zap.String("application_name", value.Name),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func validateYAML(content string) error {
	if content == "" {
		return nil
	}
	var value any
	return yaml.Unmarshal([]byte(content), &value)
}
