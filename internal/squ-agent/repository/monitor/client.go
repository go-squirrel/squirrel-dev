package monitor

import (
	"squirrel-dev/internal/squ-agent/model"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
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

func (c *Client) CreateDiskUsageMonitor(data *model.DiskUsageMonitor) error {
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

	// 删除 DiskUsageMonitor 表的过期数据
	err = c.DB.Where("collect_time < ?", beforeTime).Delete(&model.DiskUsageMonitor{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetBaseMonitorByTimeRange(since time.Time) ([]model.BaseMonitor, error) {
	var monitors []model.BaseMonitor
	err := c.DB.Where("collect_time >= ?", since).
		Order("collect_time ASC").
		Find(&monitors).Error
	return monitors, err
}

func (c *Client) GetDiskIOMonitorByTimeRange(since time.Time) ([]model.DiskIOMonitor, error) {
	var monitors []model.DiskIOMonitor
	err := c.DB.Where("collect_time >= ?", since).
		Order("collect_time ASC").
		Find(&monitors).Error
	return monitors, err
}

func (c *Client) GetNetworkMonitorByTimeRange(since time.Time) ([]model.NetworkMonitor, error) {
	var monitors []model.NetworkMonitor
	err := c.DB.Where("collect_time >= ?", since).
		Order("collect_time ASC").
		Find(&monitors).Error
	return monitors, err
}

func (c *Client) GetDiskUsageMonitorByTimeRange(since time.Time) ([]model.DiskUsageMonitor, error) {
	var monitors []model.DiskUsageMonitor
	err := c.DB.Where("collect_time >= ?", since).
		Order("collect_time ASC").
		Find(&monitors).Error
	return monitors, err
}
