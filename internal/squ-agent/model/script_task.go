package model

import "time"

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
