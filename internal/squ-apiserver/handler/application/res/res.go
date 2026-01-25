package res

type Application struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Version     string `json:"version"`
}

// 构建服务器列表响应
type ServerInfo struct {
	ServerID   uint   `json:"server_id"`
	IpAddress  string `json:"ip_address"`
	AgentPort  int    `json:"agent_port"`
	DeployedAt string `json:"deployed_at"`
}
