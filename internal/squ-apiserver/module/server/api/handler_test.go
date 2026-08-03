package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/server/api/res"
	"squirrel-dev/internal/squ-apiserver/module/server/application"
	"squirrel-dev/internal/squ-apiserver/module/server/domain"
	"squirrel-dev/pkg/jwt"
)

type repositoryStub struct {
	servers []domain.Server
	added   *domain.Server
}

func (r *repositoryStub) List(context.Context) ([]domain.Server, error) { return r.servers, nil }
func (r *repositoryStub) Get(_ context.Context, id uint) (domain.Server, error) {
	for _, server := range r.servers {
		if server.ID == id {
			return server, nil
		}
	}
	return domain.Server{}, gorm.ErrRecordNotFound
}
func (r *repositoryStub) Delete(context.Context, uint) error { return nil }
func (r *repositoryStub) Add(_ context.Context, server *domain.Server) error {
	copy := *server
	r.added = &copy
	return nil
}
func (r *repositoryStub) Update(context.Context, *domain.Server) error { return nil }
func (r *repositoryStub) GetByUUID(context.Context, string) (domain.Server, error) {
	return domain.Server{}, gorm.ErrRecordNotFound
}

type agentStub struct{}

func (agentStub) GetInfo(context.Context, string, int) (string, map[string]any) {
	return domain.StatusOnline, map[string]any{"hostname": "agent-host"}
}

type sshStub struct{ err error }

func (s sshStub) Test(context.Context, domain.Server) error { return s.err }

func TestServerHTTPContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	res.RegisterCode()
	password := "secret"
	repository := &repositoryStub{servers: []domain.Server{{
		ID: 1, Hostname: "demo", IPAddress: "192.0.2.1", AgentPort: 10750,
		SSHUsername: "root", SSHPassword: &password, SSHPort: 22, AuthType: "password",
	}}}
	service := application.NewService(repository, agentStub{}, sshStub{})
	engine := gin.New()
	RegisterRoutes(engine.Group("/api/v1"), NewHandler(service, "key"))

	assertServerRequest(t, engine, http.MethodGet, "/api/v1/server/1", "", `{"code":0,"message":"success","data":{"id":1,"hostname":"demo","ip_address":"192.0.2.1","port":10750,"ssh_username":"root","ssh_password":"secret","ssh_private_key":null,"ssh_port":22,"auth_type":"password","status":"online","server_info":{"hostname":"agent-host"}}}`)
	assertServerRequest(t, engine, http.MethodGet, "/api/v1/server/bad", "", `{"code":60021,"message":"invalid parameter"}`)
	assertServerRequest(t, engine, http.MethodPost, "/api/v1/server/check", `{}`, `{"code":60021,"message":"invalid parameter"}`)
	assertServerRequest(t, engine, http.MethodPost, "/api/v1/server", `{"ip_address":"198.51.100.2","auth_type":"password"}`, `{"code":0,"message":"success","data":"success"}`)
	if repository.added == nil || repository.added.Hostname != "198.51.100.2" || repository.added.UUID == "" {
		t.Fatalf("request mapping mismatch: %#v", repository.added)
	}
}

func TestTerminalMessageAuthenticationContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	res.RegisterCode()
	service := application.NewService(&repositoryStub{}, agentStub{}, sshStub{})
	engine := gin.New()
	group := engine.Group("/api/v1")
	handler := NewHandler(service, "websocket-key")
	RegisterRoutes(group, handler)
	RegisterTerminalRoute(group, handler)
	server := httptest.NewServer(engine)
	defer server.Close()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/server/1"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.WriteJSON(map[string]any{"type": "unexpected"}); err != nil {
		t.Fatal(err)
	}
	assertWSMessage(t, conn, "error", "expected auth message")
	_ = conn.Close()

	token, err := jwt.New("websocket-key").GenToken("demo", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.WriteJSON(map[string]any{"type": "auth", "token": token}); err != nil {
		t.Fatal(err)
	}
	assertWSMessage(t, conn, "auth_success", "authenticated")
	assertWSMessage(t, conn, "error", "server not found")
	_ = conn.Close()
}

func assertWSMessage(t *testing.T, conn *websocket.Conn, messageType, data string) {
	t.Helper()
	var value struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
	if err := conn.ReadJSON(&value); err != nil {
		t.Fatal(err)
	}
	if value.Type != messageType || value.Data != data {
		t.Fatalf("message = %#v", value)
	}
}

func assertServerRequest(t *testing.T, engine http.Handler, method, path, body, expected string) {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || recorder.Body.String() != expected {
		t.Fatalf("%s %s status=%d body=%s\nwant=%s", method, path, recorder.Code, recorder.Body.String(), expected)
	}
}
