package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/pkg/response"
	configModule "squirrel-dev/internal/squ-agent/module/config"
)

func TestConfigHTTPContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	registry := migration.NewMigrationRegistry()
	configModule.RegisterMigrations(registry)
	if err := migration.RunMigrations(db, registry); err != nil {
		t.Fatal(err)
	}

	engine := gin.New()
	configModule.RegisterHTTP(engine.Group("/api/v1"), db)

	assertRequest(t, engine, http.MethodGet, "/api/v1/config", "", `{"code":0,"message":"success","data":[{"id":1,"key":"monitor_interval","value":"300"},{"id":2,"key":"monitor_expired","value":"604800"}]}`)
	assertRequest(t, engine, http.MethodGet, "/api/v1/config/1", "", `{"code":0,"message":"success","data":{"id":1,"key":"monitor_interval","value":"300"}}`)
	assertRequest(t, engine, http.MethodGet, "/api/v1/config/not-a-number", "", `{"code":41001,"message":"parameter error"}`)
	assertRequest(t, engine, http.MethodPost, "/api/v1/config", `{"id":99,"key":"new_key","value":"new_value"}`, `{"code":0,"message":"success","data":"success"}`)
	assertRequest(t, engine, http.MethodGet, "/api/v1/config/3", "", `{"code":0,"message":"success","data":{"id":3,"key":"new_key","value":"new_value"}}`)
	assertRequest(t, engine, http.MethodDelete, "/api/v1/config/3", "", `{"code":0,"message":"success","data":"success"}`)
	assertRequest(t, engine, http.MethodGet, "/api/v1/config/3", "", `{"code":50001,"message":"sql not found"}`)
}

func TestSaveRetainsLegacyZeroValueUpdateBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	registry := migration.NewMigrationRegistry()
	configModule.RegisterMigrations(registry)
	if err := migration.RunMigrations(db, registry); err != nil {
		t.Fatal(err)
	}

	engine := gin.New()
	configModule.RegisterHTTP(engine.Group("/api/v1"), db)
	assertRequest(t, engine, http.MethodPost, "/api/v1/config", `{"key":"monitor_interval","value":""}`, `{"code":0,"message":"success","data":"success"}`)
	assertRequest(t, engine, http.MethodGet, "/api/v1/config/1", "", `{"code":0,"message":"success","data":{"id":1,"key":"monitor_interval","value":"300"}}`)
}

func assertRequest(t *testing.T, engine http.Handler, method, path, body, expected string) {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("%s %s status = %d", method, path, recorder.Code)
	}
	if recorder.Body.String() != expected {
		t.Fatalf("%s %s body = %s\nwant = %s", method, path, recorder.Body.String(), expected)
	}
}
