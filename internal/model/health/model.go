package health

import (

	"gorm.io/gorm"
)

type HealthClient struct {
	DB *gorm.DB
}

