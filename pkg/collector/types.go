package collector

import (
	"time"
)

// HostInfo 宿主机信息汇总
type HostInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	CPU       CPUInfo   `json:"cpu"`
	Memory    MemInfo   `json:"memory"`
	Disk      DiskInfo  `json:"disk"`
	OS        string    `json:"os"`
	Platform  string    `json:"platform"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	Model        string     `json:"model"`
	Cores        int        `json:"cores"`
	Usage        float64    `json:"usage"`        // 总使用率百分比
	PerCoreUsage []float64  `json:"perCoreUsage"` // 每个核心使用率
	Frequency    float64    `json:"frequency"`    // MHz
	LoadAverage  [3]float64 `json:"loadAverage"`
}

// MemInfo 内存信息
type MemInfo struct {
	Total     uint64  `json:"total"`     // 字节
	Available uint64  `json:"available"` // 字节
	Used      uint64  `json:"used"`      // 字节
	Usage     float64 `json:"usage"`     // 使用率百分比
	SwapTotal uint64  `json:"swapTotal"`
	SwapUsed  uint64  `json:"swapUsed"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Partitions []Partition `json:"partitions"`
	Total      uint64      `json:"total"`     // 字节
	Used       uint64      `json:"used"`      // 字节
	Available  uint64      `json:"available"` // 字节
	Usage      float64     `json:"usage"`     // 使用率百分比
}

// Partition 磁盘分区信息
type Partition struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mountPoint"`
	FSType     string  `json:"fsType"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Available  uint64  `json:"available"`
	Usage      float64 `json:"usage"`
}
