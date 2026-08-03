package infra

import (
	"time"

	"gorm.io/gorm"
)

type configModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Key       string
	Value     string
}

func (configModel) TableName() string {
	return "configs"
}
