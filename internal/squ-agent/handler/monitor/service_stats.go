package monitor

import (
	"sort"
	"squirrel-dev/internal/pkg/response"
	monitorres "squirrel-dev/internal/squ-agent/handler/monitor/res"
	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/collector"
)

// GetStats 获取系统统计数据
func (m *Monitor) GetStats() response.Response {
	if m.Factory == nil {
		return response.Error(model.ReturnErrCode(nil))
	}

	// 收集主机信息
	hostInfo, err := m.Factory.CollectAll()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
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

	return response.Success(stats)
}

// GetAllNetIO 获取所有网卡IO统计（汇总）
func (m *Monitor) GetAllNetIO() response.Response {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil))
	}

	ioStatsList, err := ioCollector.CollectAllNetIO()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 汇总所有网卡IO数据
	var result monitorres.AllNetIOStats
	result.Data.Name = "all"
	for _, ioStats := range ioStatsList {
		result.Data.BytesSent += ioStats.BytesSent
		result.Data.BytesRecv += ioStats.BytesRecv
		result.Data.PacketsSent += ioStats.PacketsSent
		result.Data.PacketsRecv += ioStats.PacketsRecv
		result.Data.Errin += ioStats.Errin
		result.Data.Errout += ioStats.Errout
		result.Data.Dropin += ioStats.Dropin
		result.Data.Dropout += ioStats.Dropout
		result.Ifnames = append(result.Ifnames, ioStats.Name)
	}
	sort.Strings(result.Ifnames)
	return response.Success(result)
}

// GetNetIO 获取指定网卡IO统计
func (m *Monitor) GetNetIO(interfaceName string) response.Response {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil))
	}

	ioStats, err := ioCollector.CollectNetIO(interfaceName)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
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

	return response.Success(result)
}

// GetAllDiskIO 获取所有磁盘IO统计（汇总）
func (m *Monitor) GetAllDiskIO() response.Response {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil))
	}

	ioStatsList, err := ioCollector.CollectAllDiskIO()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 汇总所有磁盘IO数据
	var allResult monitorres.AllDiskIOStats

	allResult.Data.Device = "all"
	for _, ioStats := range ioStatsList {
		allResult.Data.ReadBytes += ioStats.IOCounters.ReadBytes
		allResult.Data.WriteBytes += ioStats.IOCounters.WriteBytes
		allResult.Data.ReadCount += ioStats.IOCounters.ReadCount
		allResult.Data.WriteCount += ioStats.IOCounters.WriteCount
		allResult.Data.ReadTime += ioStats.IOCounters.ReadTime
		allResult.Data.WriteTime += ioStats.IOCounters.WriteTime
		allResult.Devices = append(allResult.Devices, ioStats.Device)
	}

	sort.Strings(allResult.Devices)
	return response.Success(allResult)
}

// GetDiskIO 获取指定磁盘IO统计
func (m *Monitor) GetDiskIO(device string) response.Response {
	ioCollector := m.Factory.GetIOCollector()
	if ioCollector == nil {
		return response.Error(model.ReturnErrCode(nil))
	}

	ioStats, err := ioCollector.CollectDiskIO(device)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
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

	return response.Success(result)
}
