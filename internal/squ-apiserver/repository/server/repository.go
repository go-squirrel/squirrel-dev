package server

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() (servers []model.Server, err error)
	Get(id uint) (servers model.Server, err error)
	Delete(id uint) (err error)
	Add(req *model.Server) (err error)
	Update(req *model.Server) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
