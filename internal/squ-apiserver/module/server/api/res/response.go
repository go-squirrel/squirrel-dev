package res

type Server struct {
	ID            uint           `json:"id"`
	Hostname      string         `json:"hostname"`
	IPAddress     string         `json:"ip_address"`
	Port          int            `json:"port"`
	SSHUsername   string         `json:"ssh_username"`
	SSHPassword   *string        `json:"ssh_password"`
	SSHPrivateKey *string        `json:"ssh_private_key"`
	SSHPort       int            `json:"ssh_port"`
	AuthType      string         `json:"auth_type"`
	Status        string         `json:"status"`
	ServerAlias   *string        `json:"server_alias,omitempty"`
	ServerInfo    map[string]any `json:"server_info"`
}

type SSHTestResult struct {
	Message   string `json:"message"`
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
	SSHPort   int    `json:"ssh_port"`
}

type AgentCheckResult struct {
	Ready      bool           `json:"ready"`
	Message    string         `json:"message"`
	ServerInfo map[string]any `json:"server_info"`
}
