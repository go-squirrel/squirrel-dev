package app_store

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func (c *Client) List() (appStores []model.AppStore, err error) {
	err = c.DB.Find(&appStores).Error
	return appStores, err
}

func (c *Client) Get(id uint) (appStore model.AppStore, err error) {
	err = c.DB.Where("id = ?", id).First(&appStore).Error
	return appStore, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&model.AppStore{}, id).Error
}

func (c *Client) Add(req *model.AppStore) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *model.AppStore) (err error) {
	return c.DB.Updates(req).Error
}
