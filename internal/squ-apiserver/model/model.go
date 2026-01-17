package model

import (
	"squirrel-dev/internal/pkg/response"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 这里可以定义接口

func ReturnErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return response.ErrSQLNotFound
	case gorm.ErrDuplicatedKey:
		return response.ErrDuplicatedKey
	}
	return response.ErrSQL
}
