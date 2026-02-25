package res

import (
	"time"
)

// Stats 系统统计信息
type Stats struct {
	Timestamp   time.Time  `json:"timestamp"`
	Hostname    string     `json:"hostname"`
	LoadAverage LoadAvg    `json:"loadAverage"`
	CPU         CPUStats   `json:"cpu"`
	Memory      MemStats   `json:"memory"`
	Disk        DiskStats  `json:"disk"`
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
	Total      uint64          `json:"total"`
	Used       uint64          `json:"used"`
	Available  uint64          `json:"available"`
	Usage      float64         `json:"usage"`
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

type AllDiskIOStats struct {
	Data    DiskIOStats `json:"data"`
	Devices []string    `json:"devices"`
}

// DiskIOStats 磁盘IO统计
type DiskIOStats struct {
	Device     string `json:"device"`
	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
	ReadCount  uint64 `json:"readCount"`
	WriteCount uint64 `json:"writeCount"`
	ReadTime   uint64 `json:"readTime"`
	WriteTime  uint64 `json:"writeTime"`
}

type AllNetIOStats struct {
	Data    NetIOStats `json:"data"`
	Ifnames []string   `json:"ifnames"`
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

// BaseMonitorResponse 基础监控响应
type BaseMonitorResponse struct {
	ID          uint      `json:"id"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	MemoryTotal uint64    `json:"memory_total"`
	MemoryUsed  uint64    `json:"memory_used"`
	DiskUsage   float64   `json:"disk_usage"`
	DiskTotal   uint64    `json:"disk_total"`
	DiskUsed    uint64    `json:"disk_used"`
	CollectTime time.Time `json:"collect_time"`
}

// DiskIOMonitorResponse 磁盘IO监控响应
type DiskIOMonitorResponse struct {
	ID             uint      `json:"id"`
	DiskName       string    `json:"disk_name"`
	ReadCount      uint64    `json:"read_count"`
	WriteCount     uint64    `json:"write_count"`
	ReadBytes      uint64    `json:"read_bytes"`
	WriteBytes     uint64    `json:"write_bytes"`
	ReadTime       uint64    `json:"read_time"`
	WriteTime      uint64    `json:"write_time"`
	IoTime         uint64    `json:"io_time"`
	WeightedIoTime uint64    `json:"weighted_io_time"`
	IopsInProgress uint64    `json:"iops_in_progress"`
	CollectTime    time.Time `json:"collect_time"`
}

// NetworkMonitorResponse 网卡流量监控响应
type NetworkMonitorResponse struct {
	ID            uint      `json:"id"`
	InterfaceName string    `json:"interface_name"`
	BytesSent     uint64    `json:"bytes_sent"`
	BytesRecv     uint64    `json:"bytes_recv"`
	PacketsSent   uint64    `json:"packets_sent"`
	PacketsRecv   uint64    `json:"packets_recv"`
	ErrIn         uint64    `json:"err_in"`
	ErrOut        uint64    `json:"err_out"`
	DropIn        uint64    `json:"drop_in"`
	DropOut       uint64    `json:"drop_out"`
	FIFOIn        uint64    `json:"fifo_in"`
	FIFOOut       uint64    `json:"fifo_out"`
	CollectTime   time.Time `json:"collect_time"`
}

// DiskUsageMonitorResponse disk usage monitor response
type DiskUsageMonitorResponse struct {
	ID          uint      `json:"id"`
	DeviceName  string    `json:"device_name"`
	MountPoint  string    `json:"mount_point"`
	FsType      string    `json:"fs_type"`
	Total       uint64    `json:"total"`
	Used        uint64    `json:"used"`
	Free        uint64    `json:"free"`
	Usage       float64   `json:"usage"`
	InodesTotal uint64    `json:"inodes_total"`
	InodesUsed  uint64    `json:"inodes_used"`
	InodesFree  uint64    `json:"inodes_free"`
	CollectTime time.Time `json:"collect_time"`
}

// PageData 分页数据
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
