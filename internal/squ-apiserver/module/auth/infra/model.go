package infra

import (
	"time"

	"gorm.io/gorm"
)

type userModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"size:50;not null;unique"`
	Password  string         `gorm:"size:100;not null"`
	Email     string         `gorm:"size:100;unique"`
	Nickname  string         `gorm:"size:50"`
	Avatar    string         `gorm:"size:255"`
	Status    int            `gorm:"default:1"`
}

func (userModel) TableName() string { return "users" }
