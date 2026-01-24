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
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}

func (c *Client) CreateBaseMonitor(data *model.BaseMonitor) error {
	return c.DB.Create(data).Error
}

func (c *Client) CreateDiskIOMonitor(data *model.DiskIOMonitor) error {
	return c.DB.Create(data).Error
}

func (c *Client) CreateNetworkMonitor(data *model.NetworkMonitor) error {
	return c.DB.Create(data).Error
}

func (c *Client) DeleteBeforeTime(beforeTime time.Time) error {
	// 删除 BaseMonitor 表的过期数据
	err := c.DB.Where("collect_time < ?", beforeTime).Delete(&model.BaseMonitor{}).Error
	if err != nil {
		return err
	}

	// 删除 DiskIOMonitor 表的过期数据
	err = c.DB.Where("collect_time < ?", beforeTime).Delete(&model.DiskIOMonitor{}).Error
	if err != nil {
		return err
	}

	// 删除 NetworkMonitor 表的过期数据
	err = c.DB.Where("collect_time < ?", beforeTime).Delete(&model.NetworkMonitor{}).Error
	if err != nil {
		return err
	}

	return nil
}
