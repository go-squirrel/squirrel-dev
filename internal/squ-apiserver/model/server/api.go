package server

import "gorm.io/gorm"

type Client interface {
}

func New(db *gorm.DB) Client {
	return &ModelClient{
		DB: db,
	}
}
