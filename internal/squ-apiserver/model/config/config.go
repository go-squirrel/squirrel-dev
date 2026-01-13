package config

func (c *Client) List() (configs []Config, err error) {
	err = c.DB.Find(&configs).Error
	return configs, err
}

func (c *Client) Get(id uint) (config Config, err error) {
	err = c.DB.Where("id = ?", id).First(&config).Error
	return config, err
}

func (c *Client) Delete(id uint) (err error) {
	return c.DB.Delete(&Config{}, id).Error
}

func (c *Client) Add(req *Config) (err error) {
	return c.DB.Create(req).Error
}

func (c *Client) Update(req *Config) (err error) {
	return c.DB.Updates(req).Error
}
