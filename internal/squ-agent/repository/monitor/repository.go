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

func (c *Client) GetBaseMonitorPage(page, pageSize int) ([]model.BaseMonitor, int64, error) {
	var monitors []model.BaseMonitor
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := c.DB.Model(&model.BaseMonitor{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = c.DB.Order("collect_time DESC").Limit(pageSize).Offset(offset).Find(&monitors).Error
	if err != nil {
		return nil, 0, err
	}

	return monitors, total, nil
}

func (c *Client) GetDiskIOMonitorPage(page, pageSize int) ([]model.DiskIOMonitor, int64, error) {
	var monitors []model.DiskIOMonitor
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := c.DB.Model(&model.DiskIOMonitor{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = c.DB.Order("collect_time DESC").Limit(pageSize).Offset(offset).Find(&monitors).Error
	if err != nil {
		return nil, 0, err
	}

	return monitors, total, nil
}

func (c *Client) GetNetworkMonitorPage(page, pageSize int) ([]model.NetworkMonitor, int64, error) {
	var monitors []model.NetworkMonitor
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := c.DB.Model(&model.NetworkMonitor{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = c.DB.Order("collect_time DESC").Limit(pageSize).Offset(offset).Find(&monitors).Error
	if err != nil {
		return nil, 0, err
	}

	return monitors, total, nil
}
