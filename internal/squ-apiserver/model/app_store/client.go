package app_store

func (c *Client) List() (appStores []AppStore, err error) {
	err = c.DB.Find(&appStores).Error
	return appStores, err
}

func (c *Client) Get(id uint) (appStore AppStore, err error) {
	err = c.DB.Where("id = ?", id).First(&appStore).Error
	return appStore, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&AppStore{}, id).Error
}

func (c *Client) Add(req *AppStore) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *AppStore) (err error) {
	return c.DB.Updates(req).Error
}
