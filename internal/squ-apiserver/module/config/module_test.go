package config

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

func TestLegacyConfigContract(t *testing.T) {
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

	assertConfig(t, engine, http.MethodGet, "/api/v1/config", "", `{"code":0,"message":"success","data":[{"id":1,"key":"registry","value":"docker.io"},{"id":2,"key":"registry_username","value":""},{"id":3,"key":"registry_password","value":""}]}`)
	assertConfig(t, engine, http.MethodGet, "/api/v1/config/bad", "", `{"code":65003,"message":"invalid config key"}`)
	assertConfig(t, engine, http.MethodPost, "/api/v1/config", `{"key":"new","value":"value"}`, `{"code":0,"message":"success","data":"success"}`)
	assertConfig(t, engine, http.MethodPost, "/api/v1/config/4", `{"key":"new","value":""}`, `{"code":0,"message":"success","data":"success"}`)
	assertConfig(t, engine, http.MethodGet, "/api/v1/config/4", "", `{"code":0,"message":"success","data":{"id":4,"key":"new","value":"value"}}`)
}

func assertConfig(t *testing.T, engine http.Handler, method, path, body, expected string) {
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
