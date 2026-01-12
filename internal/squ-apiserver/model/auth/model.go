package auth

import "gorm.io/gorm"

type Client struct {
	DB *gorm.DB
}

type User struct {
	gorm.Model
	Username string `gorm:"size:50;not null;unique"`
	Password string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;unique"`
	Nickname string `gorm:"size:50"`
	Avatar   string `gorm:"size:255"`
	Status   int    `gorm:"default:1"` // 1: 正常, 0: 禁用
}
