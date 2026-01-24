package req

type Script struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TaskID  uint   `json:"task_id"` // 任务ID，对应 ScriptResult.ID
}

type ExecuteScript struct {
	ScriptID uint   `json:"script_id"`
	ServerID uint   `json:"server_id"`
}

type ScriptResultReport struct {
	TaskID       uint   `json:"task_id"`       // 任务ID，用于定位对应的执行记录
	ScriptID     uint   `json:"script_id"`
	Output       string `json:"output"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}
