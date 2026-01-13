// The model is written using APIs to facilitate the creation of mock data during service testing.
package application

import "gorm.io/gorm"

type Repository interface {
	List() (applications []Application, err error)
	Get(id uint) (application Application, err error)
	Delete(id uint) (err error)
	Add(req *Application) (err error)
	Update(req *Application) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
