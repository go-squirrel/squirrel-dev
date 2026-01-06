package server

func (c *ModelClient) List() (servers []Server, err error) {
	err = c.DB.Find(&servers).Error
	return servers, err
}

func (c *ModelClient) Get(id uint) (servers Server, err error) {
	err = c.DB.Where("id = ?", id).First(&servers).Error
	return servers, err
}

func (c *ModelClient) Delete(id uint) (err error) {
	return c.DB.Delete(&Server{}, id).Error
}

func (c *ModelClient) Add(req *Server) (err error) {

	return c.DB.Create(req).Error
}

func (c *ModelClient) Update(req *Server) (err error) {
	return c.DB.Updates(req).Error
}
