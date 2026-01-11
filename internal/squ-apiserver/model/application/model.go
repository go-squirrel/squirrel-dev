package application

import (
	"gorm.io/gorm"
)

type AppClient struct {
	DB *gorm.DB
}
