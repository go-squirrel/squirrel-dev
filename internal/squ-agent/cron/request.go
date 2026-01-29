package cron

type ReportApplicationStatus struct {
	ApplicationID uint   `json:"application_id"`
	ServerID      uint   `json:"server_id"`
	Status        string `json:"status"` // running, stopped, error等
	DeployID      uint64 `json:"deploy_id"`
}

type ReportScriptsExcute struct {
	TaskID       uint   `json:"task_id"`
	ScriptsID    uint   `json:"script_id"`
	Output       string `json:"output"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}
