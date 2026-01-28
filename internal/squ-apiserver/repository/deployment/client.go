package deployment

import "squirrel-dev/internal/squ-apiserver/model"

func (c *Client) List(serverID uint) (dms []model.Deployment, err error) {
	if serverID > 0 {
		err = c.DB.Where("server_id = ?", serverID).Find(&dms).Error
	} else {
		err = c.DB.Find(&dms).Error
	}
	return dms, err
}

func (c *Client) Get(id uint) (dm model.Deployment, err error) {
	err = c.DB.Where("id = ?", id).First(&dm).Error
	return dm, err
}

func (c *Client) GetByServerAndApp(serverID, applicationID uint) (dm model.Deployment, err error) {
	err = c.DB.Where("server_id = ? AND application_id = ?", serverID, applicationID).First(&dm).Error
	return dm, err
}

func (c *Client) GetByDeployID(deployID uint64) (dm model.Deployment, err error) {
	err = c.DB.Where("deploy_id = ?", deployID).First(&dm).Error
	return dm, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&model.Deployment{}, id).Error
}

func (c *Client) DeleteByApplicationID(applicationID uint) (err error) {
	return c.DB.Where("application_id = ?", applicationID).Delete(&model.Deployment{}).Error
}

func (c *Client) Add(req *model.Deployment) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *model.Deployment) (err error) {
	return c.DB.Updates(req).Error
}

func (c *Client) UpdateStatus(deployID uint64, status string) (err error) {
	return c.DB.Model(&model.Deployment{}).
		Where("deploy_id = ?", deployID).
		Update("status", status).Error
}
