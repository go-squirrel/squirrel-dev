package config

import (
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

type Config struct {
	gorm.Model
	Key   string
	Value string
}
