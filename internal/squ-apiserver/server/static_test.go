package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSPAFallbackAndAPIPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterStatic(engine)
	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "api")
	})

	for _, path := range []string{"/", "/servers/terminal/1"} {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `<div id="app"></div>`) {
			t.Errorf("%s did not return embedded SPA index: status=%d body=%q", path, recorder.Code, recorder.Body.String())
		}
	}

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))
	if recorder.Code != http.StatusOK || recorder.Body.String() != "api" {
		t.Fatalf("API was intercepted by SPA middleware: status=%d body=%q", recorder.Code, recorder.Body.String())
	}
}
