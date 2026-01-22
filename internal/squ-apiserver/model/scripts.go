package model

type Script struct {
	BaseModel
	Name    string
	Content string
}

// ScriptResult 脚本执行结果
type ScriptResult struct {
	BaseModel
	ScriptID     uint   `json:"script_id" gorm:"index"`
	ServerID     uint   `json:"server_id" gorm:"index"`
	ServerIP     string `json:"server_ip"`
	AgentPort    int    `json:"agent_port"`
	Output       string `gorm:"type:text"`
	Status       string `json:"status"` // running, success, failed
	ErrorMessage string `json:"error_message" gorm:"type:text"`
}
