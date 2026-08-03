package appstore

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/appstore/infra"
)

func TestLegacyAppStoreMigrationAndValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	repository := infra.NewRepository(db)
	apps, err := repository.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 5 {
		t.Fatalf("seed count = %d", len(apps))
	}
	expectedHashes := map[string]string{
		"Nginx":         "d27b5fce63db5322f2285d8ea20f419b7fad6f71abc0266d1ea5b171de1f218a",
		"MySQL":         "39a7be1ec983fc1e786a2857376193a4c87e16f70afe4795084929c98c1d0ba2",
		"Redis":         "afffc849df549a64c2aa5ce7c0fe6ffa117545af763b55c99c4ee5694472dea6",
		"Elasticsearch": "1b3570899f7970f5996d11931c483cb83a8caf6ad8069ec4ea94242ea784f4f2",
		"Jenkins":       "4933901a87bd814b6ac23f3088a05095dc30c37636155fcb5cee1f925a916017",
	}
	for _, app := range apps {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(app.Content)))
		if hash != expectedHashes[app.Name] {
			t.Fatalf("%s template hash = %s", app.Name, hash)
		}
	}

	engine := gin.New()
	RegisterHTTP(engine.Group("/api/v1"), db)
	assertAppStore(t, engine, `{"name":"bad","type":"unknown"}`, `{"code":73004,"message":"unsupported application type"}`)
	assertAppStore(t, engine, `{"name":"bad","type":"compose","content":"not: ["}`, `{"code":73003,"message":"invalid compose content"}`)
	assertAppStore(t, engine, `{"name":"demo","type":"compose","content":"  services: {}  "}`, `{"code":0,"message":"success","data":"success"}`)
	added, err := repository.Get(t.Context(), 6)
	if err != nil {
		t.Fatal(err)
	}
	if added.Content != "services: {}" {
		t.Fatalf("trimmed content = %q", added.Content)
	}
}

func assertAppStore(t *testing.T, engine http.Handler, body, expected string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/app-store", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || recorder.Body.String() != expected {
		t.Fatalf("status=%d body=%s\nwant=%s", recorder.Code, recorder.Body.String(), expected)
	}
}
