package auth

func (c *Client) Get(username, password string) bool {
	var user User
	err := c.DB.First(&user).Where("username = ? AND password = ?", username, password).Error
	if err != nil {
		return false
	}
	return true

}
