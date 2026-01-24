// The model is written using APIs to facilitate the creation of mock data during service testing.
package monitor

import (
	"squirrel-dev/internal/squ-agent/model"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBaseMonitor(data *model.BaseMonitor) error
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

func (c *Client) DeleteBeforeTime(beforeTime time.Time) error {
	return c.DB.Unscoped().Where("collect_time < ?", beforeTime).Find(&model.BaseMonitor{}).Error
}
