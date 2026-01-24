package model

type Script struct {
	BaseModel
	Name    string
	Content string
}

// ScriptResult 脚本执行结果
type ScriptResult struct {
	BaseModel
	TaskID       uint64 `json:"task_id" gorm:"uniqueIndex"` // 唯一任务ID，用于与Agent通信
	ScriptID     uint   `json:"script_id" gorm:"index"`
	ServerID     uint   `json:"server_id" gorm:"index"`
	ServerIP     string `json:"server_ip"`
	AgentPort    int    `json:"agent_port"`
	Output       string `gorm:"type:text"`
	Status       string `json:"status"` // running, success, failed
	ErrorMessage string `json:"error_message" gorm:"type:text"`
}
