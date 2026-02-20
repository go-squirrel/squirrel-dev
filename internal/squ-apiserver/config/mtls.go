package config

// MTLS mTLS 双向认证配置
type MTLS struct {
	Enabled    bool     `mapstructure:"enabled"`
	CAFile     string   `mapstructure:"caFile"`     // CA 证书文件路径
	CertFile   string   `mapstructure:"certFile"`   // 服务端证书文件路径
	KeyFile    string   `mapstructure:"keyFile"`    // 服务端私钥文件路径
	AllowedCNs []string `mapstructure:"allowedCNs"` // 允许的客户端证书 Common Name 列表
}
