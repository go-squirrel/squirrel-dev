package res

// 构建服务器列表响应
type ServerInfo struct {
	ServerID   uint   `json:"server_id"`
	IpAddress  string `json:"ip_address"`
	AgentPort  int    `json:"agent_port"`
	DeployedAt string `json:"deployed_at"`
}
