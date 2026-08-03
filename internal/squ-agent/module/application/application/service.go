package application

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"squirrel-dev/internal/squ-agent/module/application/domain"
)

var (
	ErrDockerNotInstalled    = errors.New("docker not installed")
	ErrComposeNotFound       = errors.New("docker-compose command not found")
	ErrComposeStart          = errors.New("docker-compose start failed")
	ErrComposeCreate         = errors.New("docker-compose file creation failed")
	ErrComposeStop           = errors.New("docker-compose stop failed")
	ErrLegacyDeleteShortStop = errors.New("legacy delete returned stop result")
)

type Request struct {
	ID          uint
	Name        string
	Description string
	Type        string
	Status      string
	Content     string
	Version     string
	ServerID    uint
	DeployID    uint64
	Env         []map[string]string
}

type Result struct {
	Data any
	Err  error
}

type Service struct {
	repository domain.Repository
	configs    domain.ConfigStore
	runtime    domain.ComposeRuntime
}

func NewService(repository domain.Repository, configs domain.ConfigStore, runtime domain.ComposeRuntime) *Service {
	return &Service{repository: repository, configs: configs, runtime: runtime}
}

func (s *Service) List(ctx context.Context) ([]domain.Application, error) {
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id uint) (domain.Application, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, request Request) error {
	return s.repository.Update(ctx, &domain.Application{
		ID: request.ID, Name: request.Name, Description: request.Description,
		Type: request.Type, Status: request.Status, Content: request.Content, Version: request.Version,
	})
}

func (s *Service) Add(ctx context.Context, request Request) Result {
	if !s.runtime.DockerInstalled() {
		return Result{Err: ErrDockerNotInstalled}
	}
	if !s.runtime.ComposeAvailable() {
		return Result{Err: ErrComposeNotFound}
	}
	existing, err := s.repository.GetByDeployID(ctx, request.DeployID)
	if err == nil {
		return s.redeploy(ctx, existing, request)
	}

	app := domain.Application{
		Name: request.Name, Description: request.Description, Type: request.Type,
		Status: domain.StatusStopped, Content: request.Content, Version: request.Version,
		DeployID: request.DeployID, Env: request.Env,
	}
	if err := s.repository.Add(ctx, &app); err != nil {
		return Result{Err: err}
	}
	if err := s.configs.Save(ctx, "server_id", fmt.Sprint(request.ServerID)); err != nil {
		return Result{Err: err}
	}
	if err := s.deploy(ctx, &app); err != nil {
		return Result{Err: fmt.Errorf("%w: %v", ErrComposeCreate, err)}
	}
	return Result{Data: "Application added successfully, starting in background"}
}

func (s *Service) Start(ctx context.Context, deployID uint64) Result {
	app, err := s.repository.GetByDeployID(ctx, deployID)
	if err != nil {
		return Result{Err: err}
	}
	if app.Status != domain.StatusStopped || !s.runtime.ComposeFileExists(deployID) {
		return Result{Err: ErrComposeStart}
	}
	app.Status = domain.StatusStarting
	if err := s.repository.Update(ctx, &app); err != nil {
		return Result{Err: err}
	}
	path := s.runtime.Path(deployID)
	go func() {
		if err := s.runtime.Start(path, "docker-compose.yml"); err != nil {
			s.updateStatusToFailed(context.Background(), deployID)
		}
	}()
	return Result{Data: "Start request submitted, processing in background"}
}

func (s *Service) Stop(ctx context.Context, deployID uint64) Result {
	app, err := s.repository.GetByDeployID(ctx, deployID)
	if err != nil {
		return Result{Err: err}
	}
	if app.Status != domain.StatusRunning || !s.runtime.ComposeFileExists(deployID) {
		return Result{Err: ErrComposeStop}
	}
	path := s.runtime.Path(deployID)
	if err := s.runtime.Stop(path, "docker-compose.yml"); err != nil {
		return Result{Err: ErrComposeStop}
	}
	app.Status = domain.StatusStopped
	if err := s.repository.Update(ctx, &app); err != nil {
		return Result{Err: err}
	}
	return Result{Data: "success"}
}

func (s *Service) Delete(ctx context.Context, id uint) Result {
	app, err := s.repository.Get(ctx, id)
	if err != nil {
		return Result{Err: err}
	}
	if app.Status == domain.StatusRunning {
		stop := s.Stop(ctx, app.DeployID)
		// Preserve old's `stopRes.Code != 200` condition: successful business
		// responses use code 0, so a running app is stopped but not deleted.
		if stop.Err != nil || stop.Data != nil {
			return stop
		}
	}
	if err := s.repository.Delete(ctx, id); err != nil {
		return Result{Err: err}
	}
	return Result{Data: "success"}
}

func (s *Service) DeleteByDeployID(ctx context.Context, deployID uint64) Result {
	app, err := s.repository.GetByDeployID(ctx, deployID)
	if err != nil {
		if err.Error() == "record not found" {
			return Result{Data: "application not found, skip delete"}
		}
		return Result{Err: err}
	}
	if app.Status == domain.StatusRunning {
		_ = s.Stop(ctx, deployID)
	}
	if err := s.repository.Delete(ctx, app.ID); err != nil {
		return Result{Err: err}
	}
	return Result{Data: "success"}
}

func (s *Service) redeploy(ctx context.Context, app domain.Application, request Request) Result {
	if app.Status == domain.StatusRunning {
		_ = s.runtime.Stop(s.runtime.Path(request.DeployID), "docker-compose.yml")
	}
	err := s.repository.Transaction(ctx, func(repository domain.Repository) error {
		app.Name, app.Description, app.Type = request.Name, request.Description, request.Type
		app.Content, app.Version, app.Env = request.Content, request.Version, request.Env
		app.Status = domain.StatusStarting
		if err := repository.Update(ctx, &app); err != nil {
			return fmt.Errorf("failed to update application: %w", err)
		}
		return s.configs.Save(ctx, "server_id", fmt.Sprint(request.ServerID))
	})
	if err != nil {
		return Result{Err: err}
	}
	path, file, err := s.runtime.Prepare(app.DeployID, app.Content, app.Env)
	if err != nil {
		return Result{Err: fmt.Errorf("%w: %v", ErrComposeCreate, err)}
	}
	s.startAsync(app.Name, path, file, app.DeployID)
	return Result{Data: "Application redeployed successfully, starting in background"}
}

func (s *Service) deploy(ctx context.Context, app *domain.Application) error {
	path, file, err := s.runtime.Prepare(app.DeployID, app.Content, app.Env)
	if err != nil {
		return err
	}
	app.Status = domain.StatusStarting
	if err := s.repository.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}
	s.startAsync(app.Name, path, file, app.DeployID)
	return nil
}

func (s *Service) startAsync(_ string, path, file string, deployID uint64) {
	go func() {
		if err := s.runtime.Up(path, file); err != nil {
			zap.L().Error("Failed to start docker-compose", zap.Uint64("deploy_id", deployID), zap.Error(err))
			s.updateStatusToFailed(context.Background(), deployID)
		}
	}()
}

func (s *Service) updateStatusToFailed(ctx context.Context, deployID uint64) {
	apps, err := s.repository.List(ctx)
	if err != nil {
		return
	}
	for i := range apps {
		if apps[i].DeployID == deployID {
			apps[i].Status = domain.StatusFailed
			_ = s.repository.Update(ctx, &apps[i])
			break
		}
	}
}

func IsNotFound(err error) bool { return err == gorm.ErrRecordNotFound }
