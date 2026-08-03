package infra

import (
	"time"

	"gorm.io/gorm"
)

type serverModel struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	UUID          string         `gorm:"column:uuid;type:varchar(64);not null;unique;comment:服务器唯一标识"`
	Hostname      string         `gorm:"column:hostname;type:varchar(100);not null;unique;comment:主机名"`
	IPAddress     string         `gorm:"column:ip_address;type:varchar(45);not null;index:idx_ip_address;comment:IP地址（支持IPv6）"`
	AgentPort     int            `gorm:"column:agent_port;type:int(5);default:10750;comment:agent端口"`
	SSHUsername   string         `gorm:"column:ssh_username;type:varchar(50);default:'root';comment:SSH用户名"`
	SSHPassword   *string        `gorm:"column:ssh_password;type:varchar(255);comment:SSH密码（加密存储）"`
	SSHPrivateKey *string        `gorm:"column:ssh_private_key;type:text;comment:SSH私钥（加密存储）"`
	SSHPassphrase *string        `gorm:"column:ssh_key_passphrase;type:varchar(255);comment:密钥密码（加密存储）"`
	SSHPort       int            `gorm:"column:ssh_port;type:int(5);default:22;comment:SSH端口"`
	AuthType      string         `gorm:"column:auth_type;type:varchar(20);default:'password';comment:认证方式"`
	ServerAlias   *string        `gorm:"column:server_alias;type:varchar(100);comment:服务器别名"`
	Status        string         `gorm:"column:status;type:varchar(20);comment:状态"`
}

func (serverModel) TableName() string { return "servers" }
