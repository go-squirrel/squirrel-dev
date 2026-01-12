package server

func (c *Client) List() (servers []Server, err error) {
	err = c.DB.Find(&servers).Error
	return servers, err
}

func (c *Client) Get(id uint) (servers Server, err error) {
	err = c.DB.Where("id = ?", id).First(&servers).Error
	return servers, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&Server{}, id).Error
}

func (c *Client) Add(req *Server) (err error) {

	return c.DB.Create(req).Error
}

func (c *Client) Update(req *Server) (err error) {
	return c.DB.Updates(req).Error
}
