package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/compat/contract"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
)

func TestLegacyHealthRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()

	instance := New()
	instance.Gin = gin.New()
	instance.registerHTTPRoutes()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	instance.Gin.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	const expected = `{"code":0,"message":"success","data":"health"}`
	if recorder.Body.String() != expected {
		t.Fatalf("body = %q, want %q", recorder.Body.String(), expected)
	}
}

func TestAPIServerLegacyRouteInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := database.New("sqlite", ":memory:")
	if db == nil {
		t.Fatal("create sqlite database")
	}
	defer db.Close()
	instance := New()
	instance.Config = &config.Config{}
	instance.Gin = gin.New()
	instance.DB = db
	instance.registerHTTPRoutes()

	legacy, err := contract.LoadLegacy()
	if err != nil {
		t.Fatal(err)
	}
	service, ok := legacy.Service("squ-apiserver")
	if !ok {
		t.Fatal("apiserver contract missing")
	}
	actual := make(map[string]struct{})
	for _, route := range instance.Gin.Routes() {
		actual[route.Method+" "+route.Path] = struct{}{}
	}
	for _, route := range service.Routes {
		key := route.Method + " " + route.Path
		if _, ok := actual[key]; !ok {
			t.Errorf("legacy route missing: %s", key)
		}
	}
	if len(actual) != len(service.Routes)+1 {
		t.Fatalf("route count=%d, want %d legacy routes plus /health alias", len(actual), len(service.Routes)+1)
	}
}

func TestAgentCallbacksUseMTLSOnlyWhenEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := database.New("sqlite", ":memory:")
	if db == nil {
		t.Fatal("create sqlite database")
	}
	defer db.Close()
	instance := New()
	instance.Config = &config.Config{}
	instance.Config.MTLS.Enabled = true
	instance.Gin = gin.New()
	instance.DB = db
	instance.registerHTTPRoutes()

	for _, path := range []string{"/api/v1/deployment/report", "/api/v1/scripts/receive-result"} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, path, nil)
		instance.Gin.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("%s status=%d, want %d", path, recorder.Code, http.StatusUnauthorized)
		}
	}
}

func TestTerminalUsesMessageAuthenticationOutsideJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	db := database.New("sqlite", ":memory:")
	if db == nil {
		t.Fatal("create sqlite database")
	}
	defer db.Close()

	instance := New()
	instance.Config = &config.Config{}
	instance.Config.Auth.Jwt.SigningKey = "test-signing-key"
	instance.Gin = gin.New()
	instance.DB = db
	instance.registerHTTPRoutes()

	terminalRecorder := httptest.NewRecorder()
	terminalRequest := httptest.NewRequest(http.MethodGet, "/api/v1/ws/server/1", nil)
	instance.Gin.ServeHTTP(terminalRecorder, terminalRequest)
	if terminalRecorder.Code != http.StatusBadRequest {
		t.Fatalf("terminal status=%d, want websocket upgrade error %d; body=%q",
			terminalRecorder.Code, http.StatusBadRequest, terminalRecorder.Body.String())
	}

	serverRecorder := httptest.NewRecorder()
	serverRequest := httptest.NewRequest(http.MethodGet, "/api/v1/server", nil)
	instance.Gin.ServeHTTP(serverRecorder, serverRequest)
	if serverRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("server status=%d, want %d", serverRecorder.Code, http.StatusUnauthorized)
	}
}
