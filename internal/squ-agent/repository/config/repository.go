// The model is written using APIs to facilitate the creation of mock data during service testing.
package config

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Repository interface {
	Get(key string) (string, error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}

func (c *Client) Get(key string) (string, error) {
	var config model.Config
	err := c.DB.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}
