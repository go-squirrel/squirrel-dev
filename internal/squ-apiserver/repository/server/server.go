package server

import "squirrel-dev/internal/squ-apiserver/model"

func (c *Client) List() (servers []model.Server, err error) {
	err = c.DB.Find(&servers).Error
	return servers, err
}

func (c *Client) Get(id uint) (servers model.Server, err error) {
	err = c.DB.Where("id = ?", id).First(&servers).Error
	return servers, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&model.Server{}, id).Error
}

func (c *Client) Add(req *model.Server) (err error) {

	return c.DB.Create(req).Error
}

func (c *Client) Update(req *model.Server) (err error) {
	return c.DB.Updates(req).Error
}

func (c *Client) GetByUUID(uuid string) (server model.Server, err error) {
	return server, c.DB.Where("uuid = ?", uuid).First(&server).Error
}
