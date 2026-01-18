package req

type Application struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Content     string `json:"content"`
	Version     string `json:"version"`
}

// DeployApplication 部署应用请求
type DeployApplication struct {
	ApplicationID uint `json:"-"` // 从路径参数获取
	ServerID      uint `json:"server_id"`
}


