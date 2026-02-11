package config

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func (c *Client) List() (configs []model.Config, err error) {
	err = c.DB.Find(&configs).Error
	return configs, err
}

func (c *Client) Get(id uint) (config model.Config, err error) {
	err = c.DB.Where("id = ?", id).First(&config).Error
	return config, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&model.Config{}, id).Error
}

func (c *Client) Add(req *model.Config) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *model.Config) (err error) {
	return c.DB.Updates(req).Error
}

func (c *Client) CreateOrUpdate(req *model.Config) (err error) {
	var config model.Config
	err = c.DB.Where("key = ?", req.Key).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return c.DB.Create(req).Error
	}
	req.ID = config.ID
	return c.DB.Updates(req).Error
}

func (c *Client) GetByKey(key string) (string, error) {
	var config model.Config
	err := c.DB.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

func (c *Client) Transaction(fn func(repo Repository) error) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := &Client{DB: tx}
		return fn(txRepo)
	})
}
