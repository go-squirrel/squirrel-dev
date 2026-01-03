package auth

import "gorm.io/gorm"

type Client interface {
	Get(username, password string) bool
}

func New(db *gorm.DB) Client {
	return &ModelClient{
		DB: db,
	}
}
