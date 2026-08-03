package res

type ServerInfo struct {
	Hostname        string    `json:"hostname"`
	OS              string    `json:"os"`
	Platform        string    `json:"platform"`
	PlatformVersion string    `json:"platformVersion"`
	KernelVersion   string    `json:"kernelVersion"`
	Architecture    string    `json:"architecture"`
	Uptime          uint64    `json:"uptime"`
	UptimeStr       string    `json:"uptimeStr"`
	IPAddresses     []NetAddr `json:"ipAddresses"`
}

type NetAddr struct {
	Name string   `json:"name"`
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}
