package domain

import "context"

type NetAddr struct {
	Name string
	IPv4 []string
	IPv6 []string
}

type HostInfo struct {
	Hostname        string
	OS              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	Architecture    string
	Uptime          uint64
	UptimeStr       string
	IPAddresses     []NetAddr
}

type HostInfoCollector interface {
	CollectHostInfo(ctx context.Context) (*HostInfo, error)
}
