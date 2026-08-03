package res

type Script struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ScriptResult struct {
	ID           uint   `json:"id"`
	TaskID       uint64 `json:"task_id"`
	ScriptID     uint   `json:"script_id"`
	ServerID     uint   `json:"server_id"`
	ServerIP     string `json:"server_ip"`
	AgentPort    int    `json:"agent_port"`
	Output       string `json:"output"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	CreatedAt    string `json:"created_at"`
}
