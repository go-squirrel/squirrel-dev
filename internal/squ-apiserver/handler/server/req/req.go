package req

// Server 服务器添加/更新请求
type Server struct {
	ID            uint   `json:"id"`
	Hostname      string `json:"hostname"`
	IpAddress     string `json:"ip_address"`
	Port          int    `json:"port"`
	SshUsername   string `json:"ssh_username"`
	SshPassword   string `json:"ssh_password"`
	SshPrivateKey string `json:"ssh_private_key"`
	SshPort       int    `json:"ssh_port"`
	AuthType      string `json:"auth_type"`
	Status        string `json:"status"`
	ServerAlias   string `json:"server_alias,omitempty"`
}

type Register struct {
	Hostname  string `json:"hostname" binding:"required"`   // 主机名
	UUID      string `json:"uuid" binding:"required"`       // 服务器唯一标识
	AgentPort int    `json:"agent_port" binding:"required"` // Agent端口
	IpAddress string `json:"ip_address"`
}
