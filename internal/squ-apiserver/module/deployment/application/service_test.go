package application

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/deployment/domain"
)

type repositoryStub struct {
	deployments []domain.Deployment
	addErr      error
	deleted     bool
}

func (r *repositoryStub) List(_ context.Context, serverID uint) ([]domain.Deployment, error) {
	if serverID == 0 {
		return r.deployments, nil
	}
	var result []domain.Deployment
	for _, value := range r.deployments {
		if value.ServerID == serverID {
			result = append(result, value)
		}
	}
	return result, nil
}
func (r *repositoryStub) Get(_ context.Context, id uint) (domain.Deployment, error) {
	for _, value := range r.deployments {
		if value.ID == id {
			return value, nil
		}
	}
	return domain.Deployment{}, gorm.ErrRecordNotFound
}
func (r *repositoryStub) GetByDeployID(_ context.Context, id uint64) (domain.Deployment, error) {
	for _, value := range r.deployments {
		if value.DeployID == id {
			return value, nil
		}
	}
	return domain.Deployment{}, gorm.ErrRecordNotFound
}
func (r *repositoryStub) Delete(context.Context, uint) error { r.deleted = true; return nil }
func (r *repositoryStub) Add(_ context.Context, value *domain.Deployment) error {
	if r.addErr != nil {
		return r.addErr
	}
	value.ID = uint(len(r.deployments) + 1)
	r.deployments = append(r.deployments, *value)
	return nil
}
func (r *repositoryStub) Update(context.Context, *domain.Deployment) error { return nil }
func (r *repositoryStub) UpdateStatus(context.Context, uint64, string) error {
	return nil
}

type applicationStub map[uint]domain.Application

func (a applicationStub) Get(_ context.Context, id uint) (domain.Application, error) {
	value, ok := a[id]
	if !ok {
		return domain.Application{}, gorm.ErrRecordNotFound
	}
	return value, nil
}

type serverStub map[uint]domain.Server

func (s serverStub) Get(_ context.Context, id uint) (domain.Server, error) {
	value, ok := s[id]
	if !ok {
		return domain.Server{}, gorm.ErrRecordNotFound
	}
	return value, nil
}

type agentStub struct {
	paths []string
	err   error
}

func (a *agentStub) Post(_ context.Context, _ domain.Server, path string, _ any) error {
	a.paths = append(a.paths, path)
	return a.err
}

type idStub struct{ value uint64 }

func (i idStub) Generate() (uint64, error) { return i.value, nil }

func TestComposeConflictErrors(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		existing string
		err      error
	}{
		{"empty", "", "services: {}", ErrInvalidConfig},
		{"container", "services:\n  a:\n    container_name: same", "services:\n  b:\n    container_name: same", ErrContainerConflict},
		{"port", "services:\n  a:\n    ports: ['8080:80']", "services:\n  b:\n    ports: ['8080:81']", ErrPortConflict},
		{"volume", "services:\n  a:\n    volumes: ['data:/a']", "services:\n  b:\n    volumes: ['data:/b']", ErrVolumeConflict},
		{"network", "services: {}\nnetworks:\n  shared: {}", "services: {}\nnetworks:\n  shared: {}", ErrNetworkConflict},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := checkComposeContent(test.request, test.existing); !errors.Is(err, test.err) {
				t.Fatalf("error = %v, want %v", err, test.err)
			}
		})
	}
}

func TestDeployAndRollbackAgentPaths(t *testing.T) {
	apps := applicationStub{1: {ID: 1, Name: "demo", Type: "compose", Content: "services: {}"}}
	servers := serverStub{2: {ID: 2, IPAddress: "192.0.2.1", AgentPort: 10750}}

	repository := &repositoryStub{}
	agent := &agentStub{}
	service := NewService(repository, apps, servers, agent, idStub{value: 99})
	data, err := service.Deploy(context.Background(), 1, 2)
	if err != nil || data != "deploy success" {
		t.Fatalf("data=%q err=%v", data, err)
	}
	if len(agent.paths) != 1 || agent.paths[0] != "application" || repository.deployments[0].DeployID != 99 {
		t.Fatalf("agent paths=%#v deployments=%#v", agent.paths, repository.deployments)
	}

	repository = &repositoryStub{addErr: errors.New("insert failed")}
	agent = &agentStub{}
	service = NewService(repository, apps, servers, agent, idStub{value: 100})
	_, err = service.Deploy(context.Background(), 1, 2)
	if err == nil {
		t.Fatal("expected deployment record error")
	}
	if len(agent.paths) != 2 || agent.paths[1] != "application/delete/demo" {
		t.Fatalf("legacy rollback path mismatch: %#v", agent.paths)
	}
}

func TestStartStopUndeployPaths(t *testing.T) {
	repository := &repositoryStub{deployments: []domain.Deployment{{ID: 1, ServerID: 2, DeployID: 99}}}
	agent := &agentStub{}
	service := NewService(repository, applicationStub{}, serverStub{2: {ID: 2}}, agent, idStub{})
	if _, err := service.Stop(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Start(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Undeploy(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	expected := []string{"application/stop/99", "application/start/99", "application/delete/99"}
	for i := range expected {
		if agent.paths[i] != expected[i] {
			t.Fatalf("paths = %#v", agent.paths)
		}
	}
	if !repository.deleted {
		t.Fatal("undeploy did not soft-delete deployment")
	}
}
