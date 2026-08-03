package domain

import (
	"context"
	"time"
)

const StatsCacheKey = "monitor:stats:full"

const StatsCacheTTL = 10 * time.Second

type Cache interface {
	Get(context.Context, string) (any, error)
	Set(context.Context, string, any, time.Duration) error
}

type Collector interface {
	Stats(context.Context) (Stats, error)
	AllDiskIO(context.Context) (AllDiskIOStats, error)
	DiskIO(context.Context, string) (DiskIOStats, error)
	AllNetIO(context.Context) (AllNetIOStats, error)
	NetIO(context.Context, string) (NetIOStats, error)
}

type Repository interface {
	CreateBaseMonitor(context.Context, *BaseMonitor) error
	CreateDiskIOMonitor(context.Context, *DiskIOMonitor) error
	CreateNetworkMonitor(context.Context, *NetworkMonitor) error
	CreateDiskUsageMonitor(context.Context, *DiskUsageMonitor) error
	DeleteBeforeTime(context.Context, time.Time) error
	BaseByTimeRange(context.Context, time.Time) ([]BaseMonitor, error)
	DiskIOByTimeRange(context.Context, time.Time) ([]DiskIOMonitor, error)
	NetworkByTimeRange(context.Context, time.Time) ([]NetworkMonitor, error)
	DiskUsageByTimeRange(context.Context, time.Time) ([]DiskUsageMonitor, error)
}

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

type LoadAvg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type CPUStats struct {
	Model        string    `json:"model"`
	Cores        int       `json:"cores"`
	Usage        float64   `json:"usage"`
	PerCoreUsage []float64 `json:"perCoreUsage"`
	Frequency    float64   `json:"frequency"`
}

type MemStats struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Usage     float64 `json:"usage"`
	SwapTotal uint64  `json:"swapTotal"`
	SwapUsed  uint64  `json:"swapUsed"`
}

type DiskStats struct {
	Total      uint64          `json:"total"`
	Used       uint64          `json:"used"`
	Available  uint64          `json:"available"`
	Usage      float64         `json:"usage"`
	Partitions []DiskPartition `json:"partitions"`
}

type DiskPartition struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mountPoint"`
	FSType     string  `json:"fsType"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Available  uint64  `json:"available"`
	Usage      float64 `json:"usage"`
}

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

type BaseMonitor struct {
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

type DiskIOMonitor struct {
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

type NetworkMonitor struct {
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

type DiskUsageMonitor struct {
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
