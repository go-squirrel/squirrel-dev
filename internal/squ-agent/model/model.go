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

// ScriptExecutionTask 脚本执行任务
type ScriptExecutionTask struct {
	BaseModel
	ScriptID   uint       `json:"script_id"`
	Name       string     `json:"name"`
	Content    string     `gorm:"type:text" json:"content"`
	Status     string     `json:"status"` // pending, running, success, failed
	Output     string     `gorm:"type:text" json:"output"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	ExecutedAt *time.Time `json:"executed_at"`
	Reported   bool       `json:"reported"` // 是否已上报给 server
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
