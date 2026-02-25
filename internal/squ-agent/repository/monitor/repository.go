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
	DeleteBeforeTime(beforeTime time.Time) error
	GetBaseMonitorPage(page, pageSize int) ([]model.BaseMonitor, int64, error)
	GetDiskIOMonitorPage(page, pageSize int) ([]model.DiskIOMonitor, int64, error)
	GetNetworkMonitorPage(page, pageSize int) ([]model.NetworkMonitor, int64, error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
