// The model is written using APIs to facilitate the creation of mock data during service testing.
package health

import "gorm.io/gorm"

type Repository interface {
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
