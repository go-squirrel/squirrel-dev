package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/compat/contract"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
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

func TestAllLegacyAgentRoutesAreRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()

	instance := New()
	instance.Config = &config.Config{}
	instance.Gin = gin.New()
	instance.AgentDB = database.New("sqlite", "file:agent-routes?mode=memory&cache=shared")
	instance.AppDB = database.New("sqlite", "file:app-routes?mode=memory&cache=shared")
	instance.MonitorDB = database.New("sqlite", "file:monitor-routes?mode=memory&cache=shared")
	instance.ScriptTaskDB = database.New("sqlite", "file:script-routes?mode=memory&cache=shared")
	for _, db := range []database.DB{instance.AgentDB, instance.AppDB, instance.MonitorDB, instance.ScriptTaskDB} {
		if db == nil {
			t.Fatal("failed to create in-memory database")
		}
		defer db.Close()
	}
	instance.registerHTTPRoutes()

	registered := make(map[string]struct{})
	for _, route := range instance.Gin.Routes() {
		registered[route.Method+" "+route.Path] = struct{}{}
	}
	legacy, err := contract.LoadLegacy()
	if err != nil {
		t.Fatal(err)
	}
	service, ok := legacy.Service("squ-agent")
	if !ok {
		t.Fatal("squ-agent route contract is missing")
	}
	for _, route := range service.Routes {
		key := route.Method + " " + route.Path
		if _, ok := registered[key]; !ok {
			t.Errorf("legacy route is not registered: %s", key)
		}
	}
	if len(service.Routes) != 24 {
		t.Fatalf("legacy route count = %d", len(service.Routes))
	}
}
