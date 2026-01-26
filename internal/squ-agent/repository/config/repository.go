// The model is written using APIs to facilitate the creation of mock data during service testing.
package config

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() (configs []model.Config, err error)
	Get(id uint) (config model.Config, err error)
	Delete(id uint) (err error)
	Add(req *model.Config) (err error)
	Update(req *model.Config) (err error)
	GetByKey(key string) (string, error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
