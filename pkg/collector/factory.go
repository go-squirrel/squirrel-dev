package collector

import (
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type CollectorFactory struct {
	collectors []Collector
}

func NewCollectorFactory() *CollectorFactory {
	return &CollectorFactory{}
}

// Register 注册收集器
func (f *CollectorFactory) Register(collector Collector) {
	f.collectors = append(f.collectors, collector)
}

// CollectAll 收集所有信息
func (f *CollectorFactory) CollectAll() (*HostInfo, error) {
	hostInfo := &HostInfo{
		Timestamp: time.Now(),
	}
	hi, err := host.Info()
	if err != nil {
		return nil, err
	}
	// 获取主机名和系统信息
	hostInfo.Hostname = hi.Hostname
	hostInfo.OS = hi.OS
	hostInfo.Platform = hi.Platform

	// 收集CPU信息
	if cpuCollector := f.getCollector("cpu"); cpuCollector != nil {
		if cpu, ok := cpuCollector.(CPUCollector); ok {
			cpuInfo, err := cpu.CollectCPU()
			if err == nil {
				hostInfo.CPU = *cpuInfo
			}
		}
	}

	// 收集内存信息
	if memCollector := f.getCollector("memory"); memCollector != nil {
		if memory, ok := memCollector.(MemoryCollector); ok {
			memInfo, err := memory.CollectMemory()
			if err == nil {
				hostInfo.Memory = *memInfo
			}
		}
	}

	// 收集磁盘信息
	if diskCollector := f.getCollector("disk"); diskCollector != nil {
		if disk, ok := diskCollector.(DiskCollector); ok {
			diskInfo, err := disk.CollectDisk()
			if err == nil {
				hostInfo.Disk = *diskInfo
			}
		}
	}

	return hostInfo, nil
}

func (f *CollectorFactory) getCollector(name string) Collector {
	for _, collector := range f.collectors {
		if collector.Name() == name {
			return collector
		}
	}
	return nil
}
