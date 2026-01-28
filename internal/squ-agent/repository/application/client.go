package application

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func (c *Client) List() (applications []model.Application, err error) {
	err = c.DB.Find(&applications).Error
	return applications, err
}

func (c *Client) Get(id uint) (application model.Application, err error) {
	err = c.DB.Where("id = ?", id).First(&application).Error
	return application, err
}

func (c *Client) GetByDeployID(deployID uint64) (application model.Application, err error) {
	err = c.DB.Where("deploy_id = ?", deployID).First(&application).Error
	return application, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&model.Application{}, id).Error
}

func (c *Client) Add(req *model.Application) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *model.Application) (err error) {
	return c.DB.Updates(req).Error
}
