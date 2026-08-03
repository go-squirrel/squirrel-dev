package application

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
)

func TestLegacyApplicationYAMLContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	engine := gin.New()
	RegisterHTTP(engine.Group("/api/v1"), db)

	assertApplication(t, engine, http.MethodPost, "/api/v1/application", `{"name":"empty","content":""}`, `{"code":0,"message":"success","data":"success"}`)
	assertApplication(t, engine, http.MethodPost, "/api/v1/application", `{"name":"bad","content":"value: ["}`, `{"code":71005,"message":"invalid application configuration"}`)
	assertApplication(t, engine, http.MethodPost, "/api/v1/application", `{"name":"good","type":"compose","content":"services: {}","version":"1"}`, `{"code":0,"message":"success","data":"success"}`)
	assertApplication(t, engine, http.MethodGet, "/api/v1/application/2", "", `{"code":0,"message":"success","data":{"id":2,"name":"good","description":"","type":"compose","content":"services: {}","version":"1"}}`)
	assertApplication(t, engine, http.MethodGet, "/api/v1/application/bad", "", `{"code":71005,"message":"invalid application configuration"}`)
}

func assertApplication(t *testing.T, engine http.Handler, method, path, body, expected string) {
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
