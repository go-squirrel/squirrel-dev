// The model is written using APIs to facilitate the creation of mock data during service testing.
package application

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() (applications []model.Application, err error)
	Get(id uint) (application model.Application, err error)
	Delete(id uint) (err error)
	DeleteByName(name string) (err error)
	Add(req *model.Application) (err error)
	Update(req *model.Application) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
