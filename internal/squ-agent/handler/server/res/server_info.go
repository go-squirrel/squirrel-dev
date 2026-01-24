package res

// ServerInfo 服务器信息
type ServerInfo struct {
	Hostname        string    `json:"hostname"`
	OS              string    `json:"os"`
	Platform        string    `json:"platform"`
	PlatformVersion string    `json:"platformVersion"` // 发行版本
	KernelVersion   string    `json:"kernelVersion"`   // 内核版本
	Architecture    string    `json:"architecture"`    // 系统架构
	Uptime          uint64    `json:"uptime"`          // 运行时间（秒）
	UptimeStr       string    `json:"uptimeStr"`       // 运行时间（格式化字符串）
	IPAddresses     []NetAddr `json:"ipAddresses"`     // IP地址列表
}

// NetAddr 网络地址
type NetAddr struct {
	Name string   `json:"name"` // 网卡名称
	IPv4 []string `json:"ipv4"` // IPv4地址
	IPv6 []string `json:"ipv6"` // IPv6地址
}
