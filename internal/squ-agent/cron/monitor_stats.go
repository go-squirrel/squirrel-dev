package cron

import (
	"context"
	"time"

	monitorres "squirrel-dev/internal/squ-agent/handler/monitor/res"
	"squirrel-dev/pkg/collector"

	"go.uber.org/zap"
)

// 缓存 Key 常量
const (
	CacheKeyStatsCPU     = "monitor:stats:cpu"
	CacheKeyStatsMemory  = "monitor:stats:memory"
	CacheKeyStatsDisk    = "monitor:stats:disk"
	CacheKeyStatsProcess = "monitor:stats:process"
	CacheKeyStatsFull    = "monitor:stats:full"
)

// 缓存 TTL 常量
const (
	CPUCacheTTL     = 5 * time.Second
	MemoryCacheTTL  = 10 * time.Second
	DiskCacheTTL    = 60 * time.Second
	ProcessCacheTTL = 10 * time.Second
	FullCacheTTL    = 10 * time.Second
)

// startMonitorStats 启动监控统计缓存刷新任务
func (c *Cron) startMonitorStats() error {
	zap.L().Info("starting cron task for monitor stats cache refresh",
		zap.String("cron", "monitor_stats"))

	// 每 5 秒刷新一次缓存
	_, err := c.Cron.AddFunc("*/5 * * * * *", func() {
		c.refreshMonitorStatsCache()
	})
	if err != nil {
		zap.L().Error("failed to start cron task for monitor stats cache",
			zap.String("cron", "monitor_stats"),
			zap.Error(err))
		return err
	}

	// 启动时立即刷新一次（预热缓存）
	go c.refreshMonitorStatsCache()

	return nil
}

// refreshMonitorStatsCache 刷新监控统计缓存
func (c *Cron) refreshMonitorStatsCache() {
	ctx := context.Background()

	// 创建收集器工厂
	factory := collector.NewCollectorFactory()
	factory.Register(collector.NewCPUCollector())
	factory.Register(collector.NewMemoryCollector())
	factory.Register(collector.NewDiskCollector())
	factory.Register(collector.NewIOCollector())
	factory.Register(collector.NewProcessCollector())

	// 收集主机信息
	hostInfo, err := factory.CollectAll()
	if err != nil {
		zap.L().Error("failed to collect host info for cache",
			zap.String("cron", "monitor_stats"),
			zap.Error(err))
		return
	}

	// 分层缓存：CPU
	if hostInfo.CPU.Usage > 0 || len(hostInfo.CPU.PerCoreUsage) > 0 {
		c.Cache.Set(ctx, CacheKeyStatsCPU, hostInfo.CPU, CPUCacheTTL)
	}

	// 分层缓存：Memory
	c.Cache.Set(ctx, CacheKeyStatsMemory, hostInfo.Memory, MemoryCacheTTL)

	// 分层缓存：Disk
	c.Cache.Set(ctx, CacheKeyStatsDisk, hostInfo.Disk, DiskCacheTTL)

	// 收集进程 TopN
	var topCPU, topMemory []collector.ProcessStats
	procCollector := factory.GetProcessCollector()
	if procCollector != nil {
		topCPU, _ = procCollector.CollectTopCPU(5)
		topMemory, _ = procCollector.CollectTopMemory(5)
	}

	// 构建完整的统计数据
	stats := c.buildStatsResponse(hostInfo, topCPU, topMemory)

	// 缓存进程数据
	if len(topCPU) > 0 || len(topMemory) > 0 {
		c.Cache.Set(ctx, CacheKeyStatsProcess, ProcessCacheData{
			TopCPU:    topCPU,
			TopMemory: topMemory,
		}, ProcessCacheTTL)
	}

	// 缓存完整数据
	c.Cache.Set(ctx, CacheKeyStatsFull, stats, FullCacheTTL)

	zap.L().Debug("monitor stats cache refreshed",
		zap.String("cron", "monitor_stats"),
		zap.Float64("cpu_usage", hostInfo.CPU.Usage),
		zap.Float64("memory_usage", hostInfo.Memory.Usage))
}

// ProcessCacheData 进程缓存数据
type ProcessCacheData struct {
	TopCPU    []collector.ProcessStats
	TopMemory []collector.ProcessStats
}

// buildStatsResponse 构建统计数据响应
func (c *Cron) buildStatsResponse(hostInfo *collector.HostInfo, topCPU, topMemory []collector.ProcessStats) monitorres.Stats {
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

	return stats
}
