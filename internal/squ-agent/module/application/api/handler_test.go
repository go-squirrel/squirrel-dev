package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/application/api/res"
	appService "squirrel-dev/internal/squ-agent/module/application/application"
	"squirrel-dev/internal/squ-agent/module/application/domain"
)

type fakeRepository struct {
	apps    []domain.Application
	deleted bool
}

func (f *fakeRepository) List(context.Context) ([]domain.Application, error) { return f.apps, nil }
func (f *fakeRepository) Get(_ context.Context, id uint) (domain.Application, error) {
	for _, app := range f.apps {
		if app.ID == id {
			return app, nil
		}
	}
	return domain.Application{}, nil
}
func (f *fakeRepository) GetByDeployID(_ context.Context, deployID uint64) (domain.Application, error) {
	for _, app := range f.apps {
		if app.DeployID == deployID {
			return app, nil
		}
	}
	return domain.Application{}, nil
}
func (f *fakeRepository) Delete(context.Context, uint) error { f.deleted = true; return nil }
func (f *fakeRepository) Add(_ context.Context, app *domain.Application) error {
	f.apps = append(f.apps, *app)
	return nil
}
func (f *fakeRepository) Update(_ context.Context, app *domain.Application) error {
	for i := range f.apps {
		if f.apps[i].ID == app.ID {
			f.apps[i] = *app
		}
	}
	return nil
}
func (f *fakeRepository) Transaction(ctx context.Context, fn func(domain.Repository) error) error {
	return fn(f)
}

type fakeConfigStore struct{}

func (fakeConfigStore) Save(context.Context, string, string) error { return nil }

type fakeRuntime struct {
	docker  bool
	compose bool
	exists  bool
	stopped bool
}

func (f *fakeRuntime) DockerInstalled() bool  { return f.docker }
func (f *fakeRuntime) ComposeAvailable() bool { return f.compose }
func (f *fakeRuntime) Prepare(uint64, string, []map[string]string) (string, string, error) {
	return "/compose/1", "docker-compose.yml", nil
}
func (f *fakeRuntime) ComposeFileExists(uint64) bool { return f.exists }
func (f *fakeRuntime) Up(string, string) error       { return nil }
func (f *fakeRuntime) Start(string, string) error    { return nil }
func (f *fakeRuntime) Stop(string, string) error     { f.stopped = true; return nil }
func (f *fakeRuntime) Path(uint64) string            { return "/compose/1" }

func TestApplicationHTTPContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	res.RegisterCode()

	repository := &fakeRepository{apps: []domain.Application{{
		ID: 1, Name: "demo", Description: "desc", Type: "compose", Status: domain.StatusStopped,
		Content: "services: {}", Version: "1.2.3", DeployID: 42,
		Env: []map[string]string{{"SECRET": "not-returned"}},
	}}}
	service := appService.NewService(repository, fakeConfigStore{}, &fakeRuntime{})
	engine := gin.New()
	RegisterRoutes(engine.Group("/api/v1"), NewHandler(service))

	assertApplicationRequest(t, engine, http.MethodGet, "/api/v1/application", "", http.StatusOK, `{"code":0,"message":"success","data":[{"id":1,"name":"demo","description":"desc","type":"compose","status":"stopped","content":"services: {}","version":"1.2.3"}]}`)
	assertApplicationRequest(t, engine, http.MethodGet, "/api/v1/application/bad", "", http.StatusBadRequest, `{"code":41001,"message":"parameter error"}`)
	assertApplicationRequest(t, engine, http.MethodPost, "/api/v1/application", `{}`, http.StatusOK, `{"code":10001,"message":"docker not installed"}`)
}

func TestDeleteRunningApplicationPreservesLegacyStopWithoutDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	res.RegisterCode()

	repository := &fakeRepository{apps: []domain.Application{{
		ID: 1, Name: "demo", Status: domain.StatusRunning, DeployID: 42,
	}}}
	runtime := &fakeRuntime{exists: true}
	service := appService.NewService(repository, fakeConfigStore{}, runtime)
	engine := gin.New()
	RegisterRoutes(engine.Group("/api/v1"), NewHandler(service))

	assertApplicationRequest(t, engine, http.MethodDelete, "/api/v1/application/1", "", http.StatusOK, `{"code":0,"message":"success","data":"success"}`)
	if !runtime.stopped {
		t.Fatal("running application was not stopped")
	}
	if repository.deleted {
		t.Fatal("legacy behavior must not delete a running application after successful stop")
	}
}

func assertApplicationRequest(t *testing.T, engine http.Handler, method, path, body string, status int, expected string) {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	if recorder.Code != status {
		t.Fatalf("%s %s status = %d, want %d", method, path, recorder.Code, status)
	}
	if recorder.Body.String() != expected {
		t.Fatalf("%s %s body = %s\nwant = %s", method, path, recorder.Body.String(), expected)
	}
}
