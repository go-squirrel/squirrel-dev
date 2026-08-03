package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/pkg/jwt"
)

func TestLegacyLoginContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	conf := &config.Config{}
	conf.Auth.Jwt.SigningKey = "test-signing-key"
	conf.Auth.Jwt.Expired = 30
	engine := gin.New()
	NoAuthRegisterHTTP(engine.Group("/api/v1"), conf, db)

	assertLogin(t, engine, `{}`, `{"code":66002,"message":"invalid username or password"}`)
	assertLogin(t, engine, `{"username":"demo","password":"wrong"}`, `{"code":66002,"message":"invalid username or password"}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(`{"username":"demo","password":"squ123"}`))
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	var result struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 || result.Data.Token == "" {
		t.Fatalf("response = %s", recorder.Body.String())
	}
	claims, err := jwt.New("test-signing-key").ParseToken(result.Data.Token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Username != "demo" || claims.Issuer != "Hank" {
		t.Fatalf("claims = %#v", claims)
	}
}

func assertLogin(t *testing.T, engine http.Handler, body, expected string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || recorder.Body.String() != expected {
		t.Fatalf("status = %d body = %s, want %s", recorder.Code, recorder.Body.String(), expected)
	}
}
