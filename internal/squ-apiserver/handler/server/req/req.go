package req

type Server struct {
	ID          uint   `json:"id"`
	Hostname    string `json:"hostname"`
	IpAddress   string `json:"ip_address"`
	SshUsername string `json:"ssh_username"`
	SshPort     int    `json:"ssh_port"`
	AuthType    string `json:"auth_type"`
	Status      string `json:"status"`
	ServerAlias string `json:"server_alias,omitempty"`
}
