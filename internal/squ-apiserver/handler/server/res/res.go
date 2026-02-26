package res

type Server struct {
	ID            uint           `json:"id"`
	Hostname      string         `json:"hostname"`
	IpAddress     string         `json:"ip_address"`
	Port          int            `json:"port"`
	SshUsername   string         `json:"ssh_username"`
	SshPassword   *string        `json:"ssh_password"`
	SshPrivateKey *string        `json:"ssh_private_key"`
	SshPort       int            `json:"ssh_port"`
	AuthType      string         `json:"auth_type"`
	Status        string         `json:"status"`
	ServerAlias   *string        `json:"server_alias,omitempty"`
	ServerInfo    map[string]any `json:"server_info"`
}

// SSHTestResult represents the response of SSH connection test.
type SSHTestResult struct {
	Message   string `json:"message"`
	Hostname  string `json:"hostname"`
	IpAddress string `json:"ip_address"`
	SshPort   int    `json:"ssh_port"`
}

// AgentCheckResult 检查 Agent 是否就绪的响应
type AgentCheckResult struct {
	Ready      bool           `json:"ready"`       // Agent 是否就绪
	Message    string         `json:"message"`     // 提示信息
	ServerInfo map[string]any `json:"server_info"` // Agent 返回的服务器信息
}
