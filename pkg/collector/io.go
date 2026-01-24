package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/net"
)

// IOStats IO统计信息
type IOStats struct {
	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
	ReadCount  uint64 `json:"readCount"`
	WriteCount uint64 `json:"writeCount"`
	ReadTime   uint64 `json:"readTime"`
	WriteTime  uint64 `json:"writeTime"`
}

// DiskIOStats 磁盘IO统计
type DiskIOStats struct {
	Device        string  `json:"device"`
	IOCounters    IOStats `json:"ioCounters"`
	IOUtilization float64 `json:"ioUtilization"` // IO利用率
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

type IO struct {
	BaseCollector
	prevDiskIO map[string]disk.IOCountersStat
	prevNetIO  map[string]net.IOCountersStat
}

func NewIOCollector() *IO {
	return &IO{
		BaseCollector: BaseCollector{name: "io"},
		prevDiskIO:    make(map[string]disk.IOCountersStat),
		prevNetIO:     make(map[string]net.IOCountersStat),
	}
}

func (io *IO) Collect() (any, error) {
	return io.CollectAllDiskIO()
}

// CollectDiskIO 收集指定磁盘的IO统计
func (io *IO) CollectDiskIO(device string) (*DiskIOStats, error) {
	counters, err := disk.IOCounters(device)
	if err != nil {
		return nil, fmt.Errorf("获取磁盘IO统计失败: %w", err)
	}

	if counters == nil {
		return nil, fmt.Errorf("设备 %s 不存在", device)
	}

	stats := &DiskIOStats{
		Device: device,
		IOCounters: IOStats{
			ReadBytes:  counters[device].ReadBytes,
			WriteBytes: counters[device].WriteBytes,
			ReadCount:  counters[device].ReadCount,
			WriteCount: counters[device].WriteCount,
			ReadTime:   counters[device].ReadTime,
			WriteTime:  counters[device].WriteTime,
		},
	}

	return stats, nil
}

// CollectAllDiskIO 收集所有磁盘的IO统计
func (io *IO) CollectAllDiskIO() ([]DiskIOStats, error) {
	counters, err := disk.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("获取磁盘IO统计失败: %w", err)
	}

	var stats []DiskIOStats
	for device, counter := range counters {
		stats = append(stats, DiskIOStats{
			Device: device,
			IOCounters: IOStats{
				ReadBytes:  counter.ReadBytes,
				WriteBytes: counter.WriteBytes,
				ReadCount:  counter.ReadCount,
				WriteCount: counter.WriteCount,
				ReadTime:   counter.ReadTime,
				WriteTime:  counter.WriteTime,
			},
		})
	}

	return stats, nil
}

// CollectNetIO 收集指定网卡的IO统计
func (io *IO) CollectNetIO(interfaceName string) (*NetIOStats, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("获取网络IO统计失败: %w", err)
	}

	for _, counter := range counters {
		if counter.Name == interfaceName {
			return &NetIOStats{
				Name:        counter.Name,
				BytesSent:   counter.BytesSent,
				BytesRecv:   counter.BytesRecv,
				PacketsSent: counter.PacketsSent,
				PacketsRecv: counter.PacketsRecv,
				Errin:       counter.Errin,
				Errout:      counter.Errout,
				Dropin:      counter.Dropin,
				Dropout:     counter.Dropout,
			}, nil
		}
	}

	return nil, fmt.Errorf("网卡 %s 不存在", interfaceName)
}

// CollectAllNetIO 收集所有网卡的IO统计
func (io *IO) CollectAllNetIO() ([]NetIOStats, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("获取网络IO统计失败: %w", err)
	}

	var stats []NetIOStats
	for _, counter := range counters {
		stats = append(stats, NetIOStats{
			Name:        counter.Name,
			BytesSent:   counter.BytesSent,
			BytesRecv:   counter.BytesRecv,
			PacketsSent: counter.PacketsSent,
			PacketsRecv: counter.PacketsRecv,
			Errin:       counter.Errin,
			Errout:      counter.Errout,
			Dropin:      counter.Dropin,
			Dropout:     counter.Dropout,
		})
	}

	return stats, nil
}
