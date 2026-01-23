package monitor

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	monitorres "squirrel-dev/internal/squ-agent/handler/monitor/res"
	"squirrel-dev/internal/squ-agent/model"
	monitorRepository "squirrel-dev/internal/squ-agent/repository/monitor"
	"squirrel-dev/pkg/collector"
)

type Monitor struct {
	Config     *config.Config
	Repository monitorRepository.Client
	Factory    *collector.CollectorFactory
}

func (m *Monitor) Status() response.Response {
	return response.Success("monitor")
}

// GetStats 获取系统统计数据
func (m *Monitor) GetStats() (response.Response, error) {
	if m.Factory == nil {
		return response.Error(model.ReturnErrCode(nil)), nil
	}

	// 收集主机信息
	hostInfo, err := m.Factory.CollectAll()
	if err != nil {
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 收集进程信息
	var topCPU, topMemory []collector.ProcessStats
	procCollector := m.Factory.GetProcessCollector()
	if procCollector != nil {
		topCPU, _ = procCollector.CollectTopCPU(5)
		topMemory, _ = procCollector.CollectTopMemory(5)
	}

	// 构建响应
	stats := monitorres.Stats{
		Timestamp: hostInfo.Timestamp,
		Hostname:  hostInfo.Hostname,
		LoadAverage: monitorres.LoadAvg{
			Load1:  hostInfo.CPU.LoadAverage[0],
			Load5:  hostInfo.CPU.LoadAverage[1],
			Load15: hostInfo.CPU.LoadAverage[2],
		},
		CPU: monitorres.CPUStats{
			Model:        hostInfo.CPU.Model,
			Cores:        hostInfo.CPU.Cores,
			Usage:        hostInfo.CPU.Usage,
			PerCoreUsage: hostInfo.CPU.PerCoreUsage,
			Frequency:    hostInfo.CPU.Frequency,
		},
		Memory: monitorres.MemStats{
			Total:     hostInfo.Memory.Total,
			Available: hostInfo.Memory.Available,
			Used:      hostInfo.Memory.Used,
			Usage:     hostInfo.Memory.Usage,
			SwapTotal: hostInfo.Memory.SwapTotal,
			SwapUsed:  hostInfo.Memory.SwapUsed,
		},
		Disk: monitorres.DiskStats{
			Total:     hostInfo.Disk.Total,
			Used:      hostInfo.Disk.Used,
			Available: hostInfo.Disk.Available,
			Usage:     hostInfo.Disk.Usage,
		},
	}

	// 转换分区信息
	for _, part := range hostInfo.Disk.Partitions {
		stats.Disk.Partitions = append(stats.Disk.Partitions, monitorres.DiskPartition{
			Device:     part.Device,
			MountPoint: part.MountPoint,
			FSType:     part.FSType,
			Total:      part.Total,
			Used:       part.Used,
			Available:  part.Available,
			Usage:      part.Usage,
		})
	}

	// 转换进程信息
	for _, proc := range topCPU {
		stats.TopCPU = append(stats.TopCPU, monitorres.ProcStat{
			PID:           proc.PID,
			Name:          proc.Name,
			CPUPercent:    proc.CPUPercent,
			MemoryMB:      proc.MemoryMB,
			MemoryPercent: proc.MemoryPercent,
			Status:        proc.Status,
			CreateTime:    proc.CreateTime,
		})
	}

	for _, proc := range topMemory {
		stats.TopMemory = append(stats.TopMemory, monitorres.ProcStat{
			PID:           proc.PID,
			Name:          proc.Name,
			CPUPercent:    proc.CPUPercent,
			MemoryMB:      proc.MemoryMB,
			MemoryPercent: proc.MemoryPercent,
			Status:        proc.Status,
			CreateTime:    proc.CreateTime,
		})
	}

	return response.Success(stats), nil
}

// GetDiskIO 获取指定磁盘IO统计
func (m *Monitor) GetDiskIO(device string) (response.Response, error) {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil)), nil
	}

	ioStats, err := ioCollector.CollectDiskIO(device)
	if err != nil {
		return response.Error(model.ReturnErrCode(err)), nil
	}

	result := monitorres.DiskIOStats{
		Device:     ioStats.Device,
		ReadBytes:  ioStats.IOCounters.ReadBytes,
		WriteBytes: ioStats.IOCounters.WriteBytes,
		ReadCount:  ioStats.IOCounters.ReadCount,
		WriteCount: ioStats.IOCounters.WriteCount,
		ReadTime:   ioStats.IOCounters.ReadTime,
		WriteTime:  ioStats.IOCounters.WriteTime,
	}

	return response.Success(result), nil
}

// GetAllDiskIO 获取所有磁盘IO统计（汇总）
func (m *Monitor) GetAllDiskIO() (response.Response, error) {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil)), nil
	}

	ioStatsList, err := ioCollector.CollectAllDiskIO()
	if err != nil {
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 汇总所有磁盘IO数据
	var result monitorres.DiskIOStats
	result.Device = "all"
	for _, ioStats := range ioStatsList {
		result.ReadBytes += ioStats.IOCounters.ReadBytes
		result.WriteBytes += ioStats.IOCounters.WriteBytes
		result.ReadCount += ioStats.IOCounters.ReadCount
		result.WriteCount += ioStats.IOCounters.WriteCount
		result.ReadTime += ioStats.IOCounters.ReadTime
		result.WriteTime += ioStats.IOCounters.WriteTime
	}

	return response.Success(result), nil
}

// GetNetIO 获取指定网卡IO统计
func (m *Monitor) GetNetIO(interfaceName string) (response.Response, error) {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil)), nil
	}

	ioStats, err := ioCollector.CollectNetIO(interfaceName)
	if err != nil {
		return response.Error(model.ReturnErrCode(err)), nil
	}

	result := monitorres.NetIOStats{
		Name:        ioStats.Name,
		BytesSent:   ioStats.BytesSent,
		BytesRecv:   ioStats.BytesRecv,
		PacketsSent: ioStats.PacketsSent,
		PacketsRecv: ioStats.PacketsRecv,
		Errin:       ioStats.Errin,
		Errout:      ioStats.Errout,
		Dropin:      ioStats.Dropin,
		Dropout:     ioStats.Dropout,
	}

	return response.Success(result), nil
}

// GetAllNetIO 获取所有网卡IO统计（汇总）
func (m *Monitor) GetAllNetIO() (response.Response, error) {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil)), nil
	}

	ioStatsList, err := ioCollector.CollectAllNetIO()
	if err != nil {
		return response.Error(model.ReturnErrCode(err)), nil
	}

	// 汇总所有网卡IO数据
	var result monitorres.NetIOStats
	result.Name = "all"
	for _, ioStats := range ioStatsList {
		result.BytesSent += ioStats.BytesSent
		result.BytesRecv += ioStats.BytesRecv
		result.PacketsSent += ioStats.PacketsSent
		result.PacketsRecv += ioStats.PacketsRecv
		result.Errin += ioStats.Errin
		result.Errout += ioStats.Errout
		result.Dropin += ioStats.Dropin
		result.Dropout += ioStats.Dropout
	}

	return response.Success(result), nil
}
