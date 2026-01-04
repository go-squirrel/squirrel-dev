package server

import (
	"gorm.io/gorm"
)

type Client interface {
	List() (servers []Server, err error)
}

func New(db *gorm.DB) Client {
	return &ModelClient{
		DB: db,
	}
}
