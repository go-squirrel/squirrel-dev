package infra

import (
	"time"

	"gorm.io/gorm"
)

type scriptModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Content   string
}

func (scriptModel) TableName() string { return "scripts" }

type resultModel struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	TaskID       uint64         `json:"task_id" gorm:"uniqueIndex"`
	ScriptID     uint           `json:"script_id" gorm:"index"`
	ServerID     uint           `json:"server_id" gorm:"index"`
	ServerIP     string         `json:"server_ip"`
	AgentPort    int            `json:"agent_port"`
	Output       string         `gorm:"type:text"`
	Status       string         `json:"status"`
	ErrorMessage string         `json:"error_message" gorm:"type:text"`
}

func (resultModel) TableName() string { return "script_results" }
