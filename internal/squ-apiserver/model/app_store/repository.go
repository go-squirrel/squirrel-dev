// The model is written using APIs to facilitate the creation of mock data during service testing.
package app_store

import "gorm.io/gorm"

type Repository interface {
	List() (appStores []AppStore, err error)
	Get(id uint) (appStore AppStore, err error)
	Delete(id uint) (err error)
	Add(req *AppStore) (err error)
	Update(req *AppStore) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
