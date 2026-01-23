package res

import (
	"time"
)

// Stats 系统统计信息
type Stats struct {
	Timestamp   time.Time `json:"timestamp"`
	Hostname    string    `json:"hostname"`
	LoadAverage LoadAvg   `json:"loadAverage"`
	CPU         CPUStats  `json:"cpu"`
	Memory      MemStats  `json:"memory"`
	Disk        DiskStats `json:"disk"`
	TopCPU      []ProcStat `json:"topCPU"`
	TopMemory   []ProcStat `json:"topMemory"`
}

// LoadAvg 负载平均值
type LoadAvg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// CPUStats CPU统计
type CPUStats struct {
	Model        string    `json:"model"`
	Cores        int       `json:"cores"`
	Usage        float64   `json:"usage"`
	PerCoreUsage []float64 `json:"perCoreUsage"`
	Frequency    float64   `json:"frequency"`
}

// MemStats 内存统计
type MemStats struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Usage     float64 `json:"usage"`
	SwapTotal uint64  `json:"swapTotal"`
	SwapUsed  uint64  `json:"swapUsed"`
}

// DiskStats 磁盘统计
type DiskStats struct {
	Total      uint64       `json:"total"`
	Used       uint64       `json:"used"`
	Available  uint64       `json:"available"`
	Usage      float64      `json:"usage"`
	Partitions []DiskPartition `json:"partitions"`
}

// DiskPartition 磁盘分区
type DiskPartition struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mountPoint"`
	FSType     string  `json:"fsType"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Available  uint64  `json:"available"`
	Usage      float64 `json:"usage"`
}

// ProcStat 进程统计
type ProcStat struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryMB      float64 `json:"memoryMB"`
	MemoryPercent float32 `json:"memoryPercent"`
	Status        string  `json:"status"`
	CreateTime    int64   `json:"createTime"`
}

// DiskIOStats 磁盘IO统计
type DiskIOStats struct {
	Device   string `json:"device"`
	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
	ReadCount  uint64 `json:"readCount"`
	WriteCount uint64 `json:"writeCount"`
	ReadTime   uint64 `json:"readTime"`
	WriteTime  uint64 `json:"writeTime"`
}

// NetIOStats 网络IO统计
type NetIOStats struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
	Errin       uint64 `json:"errin"`
	Errout      uint64 `json:"errout"`
	Dropin      uint64 `json:"dropin"`
	Dropout     uint64 `json:"dropout"`
}
