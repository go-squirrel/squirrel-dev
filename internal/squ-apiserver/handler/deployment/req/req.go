package req

// DeployApplication 部署应用请求
type DeployApplication struct {
	ApplicationID uint `json:"-"` // 从路径参数获取
	ServerID      uint `json:"server_id"`
}

// ReportApplicationStatus agent汇报应用状态请求
type ReportApplicationStatus struct {
	ApplicationID uint   `json:"application_id"`
	ServerID      uint   `json:"server_id"`
	Status        string `json:"status"` // running, stopped, error等
	DeployID      uint64 `json:"deploy_id"`
}

type ApplicationAgent struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Version     string `json:"version"`
	ServerID    uint   `json:"server_id"`
	DeployID    uint64 `json:"deploy_id"`
}
