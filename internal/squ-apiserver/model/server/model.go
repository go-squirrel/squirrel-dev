package server

import "gorm.io/gorm"

type ModelClient struct {
	DB *gorm.DB
}

type Server struct {
	gorm.Model
	Hostname      string  `gorm:"column:hostname;type:varchar(100);not null;unique;comment:主机名" json:"hostname"`
	IpAddress     string  `gorm:"column:ip_address;type:varchar(45);not null;index:idx_ip_address;comment:IP地址（支持IPv6）" json:"ip_address"`
	SshUsername   string  `gorm:"column:ssh_username;type:varchar(50);default:'root';comment:SSH用户名" json:"ssh_username"`
	SshPassword   *string `gorm:"column:ssh_password;type:varchar(255);comment:SSH密码（加密存储）" json:"ssh_password,omitempty"`
	SshPrivateKey *string `gorm:"column:ssh_private_key;type:text;comment:SSH私钥（加密存储）" json:"ssh_private_key,omitempty"`
	SshPassphrase *string `gorm:"column:ssh_key_passphrase;type:varchar(255);comment:密钥密码（加密存储）" json:"ssh_key_passphrase,omitempty"`
	SshPort       int     `gorm:"column:ssh_port;type:int(5);default:22;comment:SSH端口" json:"ssh_port"`
	AuthType      string  `gorm:"column:auth_type;type:ENUM('password','key','password_key');default:'password';comment:认证方式" json:"auth_type"`
	ServerAlias   *string `gorm:"column:server_alias;type:varchar(100);comment:服务器别名" json:"server_alias,omitempty"`
	Status        string  `gorm:"column:status;type:ENUM('active','inactive','maintenance');default:'active';index:idx_status;comment:状态" json:"status"`
}

func (Server) TableName() string {
	return "servers"
}
