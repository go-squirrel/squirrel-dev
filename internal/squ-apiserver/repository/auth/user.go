package auth

import "squirrel-dev/internal/squ-apiserver/model"

func (c *Client) Get(username, password string) bool {
	var user model.User
	err := c.DB.First(&user).Where("username = ? AND password = ?", username, password).Error
	if err != nil {
		return false
	}
	return true

}
