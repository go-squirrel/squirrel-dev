package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/module/monitor/api/res"
	"squirrel-dev/internal/squ-agent/module/monitor/application"
	"squirrel-dev/internal/squ-agent/module/monitor/domain"
)

type fakeRepository struct {
	base []domain.BaseMonitor
}

func (f fakeRepository) CreateBaseMonitor(context.Context, *domain.BaseMonitor) error { return nil }
func (f fakeRepository) CreateDiskIOMonitor(context.Context, *domain.DiskIOMonitor) error {
	return nil
}
func (f fakeRepository) CreateNetworkMonitor(context.Context, *domain.NetworkMonitor) error {
	return nil
}
func (f fakeRepository) CreateDiskUsageMonitor(context.Context, *domain.DiskUsageMonitor) error {
	return nil
}
func (f fakeRepository) DeleteBeforeTime(context.Context, time.Time) error { return nil }
func (f fakeRepository) BaseByTimeRange(context.Context, time.Time) ([]domain.BaseMonitor, error) {
	return f.base, nil
}
func (f fakeRepository) DiskIOByTimeRange(context.Context, time.Time) ([]domain.DiskIOMonitor, error) {
	return nil, nil
}
func (f fakeRepository) NetworkByTimeRange(context.Context, time.Time) ([]domain.NetworkMonitor, error) {
	return nil, nil
}
func (f fakeRepository) DiskUsageByTimeRange(context.Context, time.Time) ([]domain.DiskUsageMonitor, error) {
	return nil, nil
}

type fakeCollector struct {
	stats domain.Stats
	err   error
}

func (f fakeCollector) Stats(context.Context) (domain.Stats, error) { return f.stats, f.err }
func (f fakeCollector) AllDiskIO(context.Context) (domain.AllDiskIOStats, error) {
	return domain.AllDiskIOStats{}, f.err
}
func (f fakeCollector) DiskIO(context.Context, string) (domain.DiskIOStats, error) {
	return domain.DiskIOStats{}, f.err
}
func (f fakeCollector) AllNetIO(context.Context) (domain.AllNetIOStats, error) {
	return domain.AllNetIOStats{}, f.err
}
func (f fakeCollector) NetIO(context.Context, string) (domain.NetIOStats, error) {
	return domain.NetIOStats{}, f.err
}

func TestMonitorHTTPContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	res.RegisterCode()

	collectedAt := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	service := application.NewService(nil, fakeRepository{
		base: []domain.BaseMonitor{{ID: 7, CPUUsage: 12.5, CollectTime: collectedAt}},
	}, fakeCollector{stats: domain.Stats{
		Timestamp: collectedAt,
		Hostname:  "legacy-host",
		CPU:       domain.CPUStats{PerCoreUsage: []float64{}},
	}})
	engine := gin.New()
	RegisterRoutes(engine.Group("/api/v1"), NewHandler(service))

	assertMonitorRequest(t, engine, "/api/v1/monitor/stats", `{"code":0,"message":"success","data":{"timestamp":"2026-07-25T12:00:00Z","hostname":"legacy-host","loadAverage":{"load1":0,"load5":0,"load15":0},"cpu":{"model":"","cores":0,"usage":0,"perCoreUsage":[],"frequency":0},"memory":{"total":0,"available":0,"used":0,"usage":0,"swapTotal":0,"swapUsed":0},"disk":{"total":0,"used":0,"available":0,"usage":0,"partitions":null},"topCPU":null,"topMemory":null}}`)
	assertMonitorRequest(t, engine, "/api/v1/monitor/base", `{"code":0,"message":"success","data":[{"id":7,"cpu_usage":12.5,"memory_usage":0,"memory_total":0,"memory_used":0,"disk_usage":0,"disk_total":0,"disk_used":0,"collect_time":"2026-07-25T12:00:00Z"}]}`)
	assertMonitorRequest(t, engine, "/api/v1/monitor/base?range=bad", `{"code":20002,"message":"parameter error"}`)
	assertMonitorRequest(t, engine, "/api/v1/monitor/disk", `{"code":0,"message":"success","data":null}`)
}

func TestMonitorCollectorFailureKeepsLegacySQLError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.Init()
	res.RegisterCode()
	service := application.NewService(nil, fakeRepository{}, fakeCollector{err: errors.New("collect failed")})
	engine := gin.New()
	RegisterRoutes(engine.Group("/api/v1"), NewHandler(service))
	assertMonitorRequest(t, engine, "/api/v1/monitor/stats", `{"code":50000,"message":"sql error"}`)
}

func assertMonitorRequest(t *testing.T, engine http.Handler, path, expected string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("%s status = %d", path, recorder.Code)
	}
	if recorder.Body.String() != expected {
		t.Fatalf("%s body = %s\nwant = %s", path, recorder.Body.String(), expected)
	}
}
