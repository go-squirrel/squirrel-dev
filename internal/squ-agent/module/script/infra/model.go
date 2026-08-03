package infra

import (
	"time"

	"gorm.io/gorm"
)

type taskModel struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ScriptID   uint
	TaskID     uint
	Name       string
	Content    string `gorm:"type:text"`
	Status     string
	Output     string `gorm:"type:text"`
	ErrorMsg   string `gorm:"type:text"`
	ExecutedAt *time.Time
	Reported   bool
}

func (taskModel) TableName() string { return "script_execution_tasks" }
