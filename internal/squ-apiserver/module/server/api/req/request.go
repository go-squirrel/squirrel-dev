package req

type Server struct {
	ID            uint   `json:"id"`
	Hostname      string `json:"hostname"`
	IPAddress     string `json:"ip_address"`
	Port          int    `json:"port"`
	SSHUsername   string `json:"ssh_username"`
	SSHPassword   string `json:"ssh_password"`
	SSHPrivateKey string `json:"ssh_private_key"`
	SSHPort       int    `json:"ssh_port"`
	AuthType      string `json:"auth_type"`
	Status        string `json:"status"`
	ServerAlias   string `json:"server_alias,omitempty"`
}

type CheckAgent struct {
	IPAddress string `json:"ip_address" binding:"required"`
	Port      int    `json:"port" binding:"required"`
}
