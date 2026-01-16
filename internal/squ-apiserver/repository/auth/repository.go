package auth

import "gorm.io/gorm"

type Repository interface {
	Get(username, password string) bool
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
