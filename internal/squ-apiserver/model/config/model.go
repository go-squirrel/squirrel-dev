package config

import (
	"gorm.io/gorm"
)

type ConfigClient struct {
	DB *gorm.DB
}
