package application_server

import "squirrel-dev/internal/squ-apiserver/model"

func (c *Client) List(serverID uint) (applicationServers []model.ApplicationServer, err error) {
	if serverID > 0 {
		err = c.DB.Where("server_id = ?", serverID).Find(&applicationServers).Error
	} else {
		err = c.DB.Find(&applicationServers).Error
	}
	return applicationServers, err
}

func (c *Client) Get(id uint) (applicationServer model.ApplicationServer, err error) {
	err = c.DB.Where("id = ?", id).First(&applicationServer).Error
	return applicationServer, err
}

func (c *Client) GetByServerAndApp(serverID, applicationID uint) (applicationServer model.ApplicationServer, err error) {
	err = c.DB.Where("server_id = ? AND application_id = ?", serverID, applicationID).First(&applicationServer).Error
	return applicationServer, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&model.ApplicationServer{}, id).Error
}

func (c *Client) DeleteByApplicationID(applicationID uint) (err error) {
	return c.DB.Where("application_id = ?", applicationID).Delete(&model.ApplicationServer{}).Error
}

func (c *Client) Add(req *model.ApplicationServer) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *model.ApplicationServer) (err error) {
	return c.DB.Updates(req).Error
}

func (c *Client) UpdateStatus(serverID, applicationID uint, status string) (err error) {
	return c.DB.Model(&model.ApplicationServer{}).
		Where("server_id = ? AND application_id = ?", serverID, applicationID).
		Update("status", status).Error
}

