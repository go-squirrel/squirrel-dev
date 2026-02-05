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
