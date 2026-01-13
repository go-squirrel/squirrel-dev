package application

import (
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

type Application struct {
	gorm.Model
	Name        string
	Description string
	Type        string // compose and manifest
	Status      string
	Content     string
	Version     string
}
