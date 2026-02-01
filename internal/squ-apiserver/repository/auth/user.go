package auth

import (
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/hash"
)

func (c *Client) Get(username, password string) bool {
	var user model.User
	err := c.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return false
	}
	err = hash.ComparePassword(user.Password, password)
	return err == nil
}
