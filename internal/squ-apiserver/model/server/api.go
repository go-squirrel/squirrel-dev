package server

import (
	"gorm.io/gorm"
)

type Client interface {
	List() (servers []Server, err error)
	Get(id uint) (servers Server, err error)
	Delete(id uint) (err error)
	Add(req *Server) (err error)
	Update(req *Server) (err error)
}

func New(db *gorm.DB) Client {
	return &ModelClient{
		DB: db,
	}
}
