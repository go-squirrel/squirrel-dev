package infra

import (
	"context"
	"fmt"
	"sort"

	"squirrel-dev/internal/squ-agent/module/monitor/domain"
	"squirrel-dev/pkg/collector"
)

type MetricsCollector struct {
	factory *collector.CollectorFactory
}

func NewMetricsCollector(factory *collector.CollectorFactory) *MetricsCollector {
	return &MetricsCollector{factory: factory}
}

func NewDefaultCollector() *MetricsCollector {
	factory := collector.NewCollectorFactory()
	factory.Register(collector.NewCPUCollector())
	factory.Register(collector.NewMemoryCollector())
	factory.Register(collector.NewDiskCollector())
	factory.Register(collector.NewIOCollector())
	factory.Register(collector.NewProcessCollector())
	return NewMetricsCollector(factory)
}

func (c *MetricsCollector) Stats(context.Context) (domain.Stats, error) {
	if c.factory == nil {
		return domain.Stats{}, fmt.Errorf("collector factory is nil")
	}
	hostInfo, err := c.factory.CollectAll()
	if err != nil {
		return domain.Stats{}, err
	}
	var topCPU, topMemory []collector.ProcessStats
	if processCollector := c.factory.GetProcessCollector(); processCollector != nil {
		topCPU, _ = processCollector.CollectTopCPU(5)
		topMemory, _ = processCollector.CollectTopMemory(5)
	}

	return BuildStats(hostInfo, topCPU, topMemory), nil
}

func BuildStats(hostInfo *collector.HostInfo, topCPU, topMemory []collector.ProcessStats) domain.Stats {
	stats := domain.Stats{
		Timestamp: hostInfo.Timestamp,
		Hostname:  hostInfo.Hostname,
		LoadAverage: domain.LoadAvg{
			Load1: hostInfo.CPU.LoadAverage[0], Load5: hostInfo.CPU.LoadAverage[1], Load15: hostInfo.CPU.LoadAverage[2],
		},
		CPU: domain.CPUStats{
			Model: hostInfo.CPU.Model, Cores: hostInfo.CPU.Cores, Usage: hostInfo.CPU.Usage,
			PerCoreUsage: hostInfo.CPU.PerCoreUsage, Frequency: hostInfo.CPU.Frequency,
		},
		Memory: domain.MemStats{
			Total: hostInfo.Memory.Total, Available: hostInfo.Memory.Available, Used: hostInfo.Memory.Used,
			Usage: hostInfo.Memory.Usage, SwapTotal: hostInfo.Memory.SwapTotal, SwapUsed: hostInfo.Memory.SwapUsed,
		},
		Disk: domain.DiskStats{
			Total: hostInfo.Disk.Total, Used: hostInfo.Disk.Used,
			Available: hostInfo.Disk.Available, Usage: hostInfo.Disk.Usage,
		},
	}
	for _, item := range hostInfo.Disk.Partitions {
		stats.Disk.Partitions = append(stats.Disk.Partitions, domain.DiskPartition{
			Device: item.Device, MountPoint: item.MountPoint, FSType: item.FSType,
			Total: item.Total, Used: item.Used, Available: item.Available, Usage: item.Usage,
		})
	}
	stats.TopCPU = processStats(topCPU)
	stats.TopMemory = processStats(topMemory)
	return stats
}

func (c *MetricsCollector) AllDiskIO(context.Context) (domain.AllDiskIOStats, error) {
	ioCollector, err := c.ioCollector()
	if err != nil {
		return domain.AllDiskIOStats{}, err
	}
	items, err := ioCollector.CollectAllDiskIO()
	if err != nil {
		return domain.AllDiskIOStats{}, err
	}
	var result domain.AllDiskIOStats
	result.Data.Device = "all"
	for _, item := range items {
		result.Data.ReadBytes += item.IOCounters.ReadBytes
		result.Data.WriteBytes += item.IOCounters.WriteBytes
		result.Data.ReadCount += item.IOCounters.ReadCount
		result.Data.WriteCount += item.IOCounters.WriteCount
		result.Data.ReadTime += item.IOCounters.ReadTime
		result.Data.WriteTime += item.IOCounters.WriteTime
		result.Devices = append(result.Devices, item.Device)
	}
	sort.Strings(result.Devices)
	return result, nil
}

func (c *MetricsCollector) DiskIO(_ context.Context, device string) (domain.DiskIOStats, error) {
	ioCollector, err := c.ioCollector()
	if err != nil {
		return domain.DiskIOStats{}, err
	}
	item, err := ioCollector.CollectDiskIO(device)
	if err != nil {
		return domain.DiskIOStats{}, err
	}
	return domain.DiskIOStats{
		Device: item.Device, ReadBytes: item.IOCounters.ReadBytes, WriteBytes: item.IOCounters.WriteBytes,
		ReadCount: item.IOCounters.ReadCount, WriteCount: item.IOCounters.WriteCount,
		ReadTime: item.IOCounters.ReadTime, WriteTime: item.IOCounters.WriteTime,
	}, nil
}

func (c *MetricsCollector) AllNetIO(context.Context) (domain.AllNetIOStats, error) {
	ioCollector, err := c.ioCollector()
	if err != nil {
		return domain.AllNetIOStats{}, err
	}
	items, err := ioCollector.CollectAllNetIO()
	if err != nil {
		return domain.AllNetIOStats{}, err
	}
	var result domain.AllNetIOStats
	result.Data.Name = "all"
	for _, item := range items {
		result.Data.BytesSent += item.BytesSent
		result.Data.BytesRecv += item.BytesRecv
		result.Data.PacketsSent += item.PacketsSent
		result.Data.PacketsRecv += item.PacketsRecv
		result.Data.Errin += item.Errin
		result.Data.Errout += item.Errout
		result.Data.Dropin += item.Dropin
		result.Data.Dropout += item.Dropout
		result.Ifnames = append(result.Ifnames, item.Name)
	}
	sort.Strings(result.Ifnames)
	return result, nil
}

func (c *MetricsCollector) NetIO(_ context.Context, interfaceName string) (domain.NetIOStats, error) {
	ioCollector, err := c.ioCollector()
	if err != nil {
		return domain.NetIOStats{}, err
	}
	item, err := ioCollector.CollectNetIO(interfaceName)
	if err != nil {
		return domain.NetIOStats{}, err
	}
	return domain.NetIOStats{
		Name: item.Name, BytesSent: item.BytesSent, BytesRecv: item.BytesRecv,
		PacketsSent: item.PacketsSent, PacketsRecv: item.PacketsRecv,
		Errin: item.Errin, Errout: item.Errout, Dropin: item.Dropin, Dropout: item.Dropout,
	}, nil
}

func (c *MetricsCollector) ioCollector() (collector.IOCollector, error) {
	if c.factory == nil || c.factory.GetIOCollector() == nil {
		return nil, fmt.Errorf("io collector is nil")
	}
	return c.factory.GetIOCollector(), nil
}

func processStats(items []collector.ProcessStats) []domain.ProcStat {
	var result []domain.ProcStat
	for _, item := range items {
		result = append(result, domain.ProcStat{
			PID: item.PID, Name: item.Name, CPUPercent: item.CPUPercent, MemoryMB: item.MemoryMB,
			MemoryPercent: item.MemoryPercent, Status: item.Status, CreateTime: item.CreateTime,
		})
	}
	return result
}
