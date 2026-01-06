package model

import (
	"squirrel-dev/internal/pkg/response"

	"gorm.io/gorm"
)

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
