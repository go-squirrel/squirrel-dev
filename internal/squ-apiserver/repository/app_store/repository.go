// The model is written using APIs to facilitate the creation of mock data during service testing.
package app_store

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

type Repository interface {
	List() (appStores []model.AppStore, err error)
	Get(id uint) (appStore model.AppStore, err error)
	Delete(id uint) (err error)
	Add(req *model.AppStore) (err error)
	Update(req *model.AppStore) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
