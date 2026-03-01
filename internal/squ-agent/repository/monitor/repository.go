// The model is written using APIs to facilitate the creation of mock data during service testing.
package monitor

import (
	"squirrel-dev/internal/squ-agent/model"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBaseMonitor(data *model.BaseMonitor) error
	CreateDiskIOMonitor(data *model.DiskIOMonitor) error
	CreateNetworkMonitor(data *model.NetworkMonitor) error
	CreateDiskUsageMonitor(data *model.DiskUsageMonitor) error
	DeleteBeforeTime(beforeTime time.Time) error
	GetBaseMonitorByTimeRange(since time.Time) ([]model.BaseMonitor, error)
	GetDiskIOMonitorByTimeRange(since time.Time) ([]model.DiskIOMonitor, error)
	GetNetworkMonitorByTimeRange(since time.Time) ([]model.NetworkMonitor, error)
	GetDiskUsageMonitorByTimeRange(since time.Time) ([]model.DiskUsageMonitor, error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
