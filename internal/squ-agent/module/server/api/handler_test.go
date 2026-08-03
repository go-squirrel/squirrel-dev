package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/server/application"
	"squirrel-dev/internal/squ-agent/module/server/domain"
)

type fakeCollector struct {
	value *domain.HostInfo
	err   error
}

func (f fakeCollector) CollectHostInfo(context.Context) (*domain.HostInfo, error) {
	return f.value, f.err
}

func TestInfoResponseContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()

	service := application.NewService(fakeCollector{value: &domain.HostInfo{
		Hostname:        "legacy-host",
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "24.04",
		KernelVersion:   "6.8",
		Architecture:    "amd64",
		Uptime:          90,
		UptimeStr:       "1分钟",
		IPAddresses: []domain.NetAddr{{
			Name: "eth0",
			IPv4: []string{"192.0.2.10"},
			IPv6: []string{"2001:db8::10"},
		}},
	}})
	engine := gin.New()
	engine.GET("/api/v1/server/info", NewHandler(service).Info)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/server/info", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	const expected = `{"code":0,"message":"success","data":{"hostname":"legacy-host","os":"linux","platform":"ubuntu","platformVersion":"24.04","kernelVersion":"6.8","architecture":"amd64","uptime":90,"uptimeStr":"1分钟","ipAddresses":[{"name":"eth0","ipv4":["192.0.2.10"],"ipv6":["2001:db8::10"]}]}}`
	if recorder.Body.String() != expected {
		t.Fatalf("body = %s\nwant = %s", recorder.Body.String(), expected)
	}
}

func TestInfoFailureKeepsLegacyError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()

	service := application.NewService(fakeCollector{err: errors.New("collect failed")})
	engine := gin.New()
	engine.GET("/api/v1/server/info", NewHandler(service).Info)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/server/info", nil))

	const expected = `{"code":50000,"message":"sql error"}`
	if recorder.Body.String() != expected {
		t.Fatalf("body = %s, want %s", recorder.Body.String(), expected)
	}
}
