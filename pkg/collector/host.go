package collector

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

// HostInfoDetail 主机详细信息
type HostInfoDetail struct {
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

type Host struct {
	BaseCollector
}

func NewHostCollector() *Host {
	return &Host{
		BaseCollector: BaseCollector{name: "host"},
	}
}

func (h *Host) Collect() (any, error) {
	return h.CollectHostInfo()
}

// CollectHostInfo 收集主机详细信息
func (h *Host) CollectHostInfo() (*HostInfoDetail, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("获取主机信息失败: %w", err)
	}

	uptime, err := host.Uptime()
	if err != nil {
		uptime = 0
	}

	// 获取IP地址（过滤虚拟网卡）
	ips := h.getIPAddresses()

	hostInfo := &HostInfoDetail{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		Architecture:    info.KernelArch,
		Uptime:          uptime,
		UptimeStr:       formatUptime(uptime),
		IPAddresses:     ips,
	}

	return hostInfo, nil
}

// getIPAddresses 获取主机的IP地址（排除虚拟网卡）
func (h *Host) getIPAddresses() []NetAddr {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var netAddrs []NetAddr

	for _, iface := range interfaces {
		// 跳过虚拟网卡
		if h.isVirtualInterface(iface.Name) {
			continue
		}

		// 跳过没有IP地址的网卡
		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}

		netAddr := NetAddr{
			Name: iface.Name,
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP

			// 跳过loopback地址
			if ip.IsLoopback() {
				continue
			}

			if ip.To4() != nil {
				// IPv4
				netAddr.IPv4 = append(netAddr.IPv4, ip.String())
			} else {
				// IPv6
				netAddr.IPv6 = append(netAddr.IPv6, ip.String())
			}
		}

		// 只添加有有效IP地址的网卡
		if len(netAddr.IPv4) > 0 || len(netAddr.IPv6) > 0 {
			netAddrs = append(netAddrs, netAddr)
		}
	}

	return netAddrs
}

// isVirtualInterface 判断是否为虚拟网卡
func (h *Host) isVirtualInterface(name string) bool {
	virtualPrefixes := []string{
		"docker", "k8s", "kube", "flannel", "cni", "calico",
		"veth", "virbr", "tun", "tap", "vif", "vni",
		"br-", "ovs", "vxlan", "geneve", "gre",
		"ip_vti", "ip6tnl", "sit", "ip6gre",
	}

	nameLower := strings.ToLower(name)
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(nameLower, prefix) {
			return true
		}
	}

	return false
}

// formatUptime 格式化运行时间
func formatUptime(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	}
	return fmt.Sprintf("%d分钟", minutes)
}
