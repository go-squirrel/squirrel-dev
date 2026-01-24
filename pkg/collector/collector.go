package collector

// Collector 收集器接口
type Collector interface {
	Name() string
	Collect() (any, error)
}

type BaseCollector struct {
	name string
}

func (b *BaseCollector) Name() string {
	return b.name
}

// CPUCollector CPU收集器接口
type CPUCollector interface {
	CollectCPU() (*CPUInfo, error)
}

// MemoryCollector 内存收集器接口
type MemoryCollector interface {
	CollectMemory() (*MemInfo, error)
}

// DiskCollector 磁盘收集器接口
type DiskCollector interface {
	CollectDisk() (*DiskInfo, error)
}

// IOCollector IO收集器接口
type IOCollector interface {
	CollectDiskIO(device string) (*DiskIOStats, error)
	CollectAllDiskIO() ([]DiskIOStats, error)
	CollectNetIO(interfaceName string) (*NetIOStats, error)
	CollectAllNetIO() ([]NetIOStats, error)
}

// ProcessCollector 进程收集器接口
type ProcessCollector interface {
	CollectTopCPU(limit int) ([]ProcessStats, error)
	CollectTopMemory(limit int) ([]ProcessStats, error)
	CollectAllProcesses() ([]ProcessStats, error)
}

// HostCollector 主机收集器接口
type HostCollector interface {
	CollectHostInfo() (*HostInfoDetail, error)
}
