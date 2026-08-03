package application

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/appstore/domain"
	"squirrel-dev/pkg/compose"
)

type Service struct{ repository domain.Repository }

func NewService(repository domain.Repository) *Service { return &Service{repository: repository} }

func (s *Service) List(ctx context.Context) ([]domain.App, error) {
	apps, err := s.repository.List(ctx)
	if err != nil {
		zap.L().Error("failed to list app stores", zap.Error(err))
		return nil, err
	}
	return apps, nil
}

func (s *Service) Get(ctx context.Context, id uint) (domain.App, error) {
	app, err := s.repository.Get(ctx, id)
	if err != nil {
		zap.L().Error("failed to get app store", zap.Uint("app_store_id", id), zap.Error(err))
		return domain.App{}, err
	}
	return app, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		zap.L().Error("failed to delete app store", zap.Uint("app_store_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Add(ctx context.Context, app domain.App) error {
	if err := validate(&app); err != nil {
		zap.L().Warn("invalid app store configuration",
			zap.String("app_store_name", app.Name),
			zap.String("application_type", app.Type),
			zap.Error(err),
		)
		return err
	}
	if err := s.repository.Add(ctx, &app); err != nil {
		zap.L().Error("failed to add app store", zap.String("app_store_name", app.Name), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, app domain.App) error {
	if err := validate(&app); err != nil {
		zap.L().Warn("invalid app store configuration",
			zap.Uint("app_store_id", app.ID),
			zap.String("app_store_name", app.Name),
			zap.String("application_type", app.Type),
			zap.Error(err),
		)
		return err
	}
	if err := s.repository.Update(ctx, &app); err != nil {
		zap.L().Error("failed to update app store",
			zap.Uint("app_store_id", app.ID),
			zap.String("app_store_name", app.Name),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func validate(app *domain.App) error {
	if !compose.IsValidAppType(app.Type) {
		return ErrUnsupportedType
	}
	if app.Type == "compose" {
		app.Content = compose.TrimSpaceContent(app.Content)
		if err := compose.ValidateContent(app.Name, app.Content); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidCompose, err)
		}
	}
	return nil
}
