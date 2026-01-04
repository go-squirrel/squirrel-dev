package server

func (c *ModelClient) List() (servers []Server, err error) {
	err = c.DB.Find(&servers).Error
	return servers, err
}
