package res

// ApplicationInfo 应用信息
type ApplicationInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Version     string `json:"version"`
}

// ServerInfo 服务器信息
type ServerInfo struct {
	ID        uint   `json:"id"`
	IpAddress string `json:"ip_address"`
	AgentPort int    `json:"agent_port"`
}

// Deployment 部署信息
type Deployment struct {
	ID          uint            `json:"id"`
	DeployID    uint64          `json:"deploy_id"`
	Application ApplicationInfo `json:"application"`
	Server      ServerInfo      `json:"server"`
	Status      string          `json:"status"`
	DeployedAt  string          `json:"deployed_at"`
	Content     string          `json:"content"`
}
