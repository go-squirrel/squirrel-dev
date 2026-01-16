package config

import "squirrel-dev/internal/squ-apiserver/model"

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
