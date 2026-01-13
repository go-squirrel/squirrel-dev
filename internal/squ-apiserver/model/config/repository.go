// The model is written using APIs to facilitate the creation of mock data during service testing.
package config

import "gorm.io/gorm"

type Repository interface {
	List() (configs []Config, err error)
	Get(id uint) (config Config, err error)
	Delete(id uint) (err error)
	Add(req *Config) (err error)
	Update(req *Config) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
