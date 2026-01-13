package application

func (c *Client) List() (applications []Application, err error) {
	err = c.DB.Find(&applications).Error
	return applications, err
}

func (c *Client) Get(id uint) (application Application, err error) {
	err = c.DB.Where("id = ?", id).First(&application).Error
	return application, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&Application{}, id).Error
}

func (c *Client) Add(req *Application) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *Application) (err error) {
	return c.DB.Updates(req).Error
}
