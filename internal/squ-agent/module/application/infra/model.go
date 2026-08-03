package infra

import (
	"time"

	"gorm.io/gorm"
)

type applicationModel struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Description string
	Type        string
	OldStatus   string
	Status      string
	Content     string
	Version     string
	DeployID    uint64
	Env         []map[string]string `gorm:"type:json;serializer:json"`
}

func (applicationModel) TableName() string { return "applications" }
