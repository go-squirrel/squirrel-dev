// The model is written using APIs to facilitate the creation of mock data during service testing.
package health

import "gorm.io/gorm"

type Client interface {
}

func New(db *gorm.DB) Client {
	return &HealthClient{
		DB: db,
	}
}
