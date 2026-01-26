package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"squirrel-dev/internal/squ-agent/model"
	configRepo "squirrel-dev/internal/squ-agent/repository/config"
	monitorRepo "squirrel-dev/internal/squ-agent/repository/monitor"
	"squirrel-dev/pkg/collector"

	"go.uber.org/zap"
)

// startMonitor 启动主机监控定时任务
func (c *Cron) startMonitor(configDBRepo configRepo.Repository, monitorDBRepo monitorRepo.Repository) error {
	// 从配置中获取监控间隔时间
	intervalSeconds, err := getMonitorInterval(configDBRepo)
	if err != nil {
		zap.L().Error("获取监控间隔配置失败，使用默认值300秒", zap.Error(err))
		intervalSeconds = 300
	}

	zap.L().Info("启动主机监控定时任务", zap.Int("interval_seconds", intervalSeconds))

	// 添加定时任务，按照配置的间隔时间执行
	intervalminites := intervalSeconds / 60
	cronTime := fmt.Sprintf("0 */%d * * * *", intervalminites)
	_, err = c.Cron.AddFunc(cronTime, func() {
		c.collectAndSaveMonitorData(configDBRepo, monitorDBRepo)
	})
	if err != nil {
		zap.L().Error("启动主机监控定时任务失败", zap.Error(err))
		return err
	}

	return nil
}

// collectAndSaveMonitorData 收集并保存监控数据
func (c *Cron) collectAndSaveMonitorData(configRepo configRepo.Repository, monitorRepo monitorRepo.Repository) {
	zap.L().Info("开始收集监控数据")

	// 收集CPU信息
	cpuCollector := collector.NewCPUCollector()
	cpuInfo, err := cpuCollector.CollectCPU()
	if err != nil {
		zap.L().Error("收集CPU信息失败", zap.Error(err))
		return
	}

	// 收集内存信息
	memCollector := collector.NewMemoryCollector()
	memInfo, err := memCollector.CollectMemory()
	if err != nil {
		zap.L().Error("收集内存信息失败", zap.Error(err))
		return
	}

	// 收集磁盘信息
	diskCollector := collector.NewDiskCollector()
	diskInfo, err := diskCollector.CollectDisk()
	if err != nil {
		zap.L().Error("收集磁盘信息失败", zap.Error(err))
		return
	}

	// 收集磁盘IO信息
	ioCollector := collector.NewIOCollector()
	diskIOStats, err := ioCollector.CollectAllDiskIO()
	if err != nil {
		zap.L().Error("收集磁盘IO信息失败", zap.Error(err))
		return
	}

	// 收集网卡流量信息
	netIOStats, err := ioCollector.CollectAllNetIO()
	if err != nil {
		zap.L().Error("收集网卡流量信息失败", zap.Error(err))
		return
	}

	// 构建基础监控数据
	baseMonitor := &model.BaseMonitor{
		CPUUsage:    cpuInfo.Usage,
		MemoryUsage: memInfo.Usage,
		MemoryTotal: memInfo.Total,
		MemoryUsed:  memInfo.Used,
		DiskUsage:   diskInfo.Usage,
		DiskTotal:   diskInfo.Total,
		DiskUsed:    diskInfo.Used,
		CollectTime: time.Now(),
	}

	// 保存监控数据
	err = monitorRepo.CreateBaseMonitor(baseMonitor)
	if err != nil {
		zap.L().Error("保存监控数据失败", zap.Error(err))
		return
	}

	zap.L().Info("监控数据保存成功",
		zap.Float64("cpu_usage", cpuInfo.Usage),
		zap.Float64("memory_usage", memInfo.Usage),
		zap.Float64("disk_usage", diskInfo.Usage),
	)

	// 保存磁盘IO监控数据
	collectTime := time.Now()
	for _, diskIO := range diskIOStats {
		// 跳过虚拟设备
		if c.shouldSkipDisk(diskIO.Device) {
			continue
		}

		diskIOMonitor := &model.DiskIOMonitor{
			DiskName:    diskIO.Device,
			ReadCount:   diskIO.IOCounters.ReadCount,
			WriteCount:  diskIO.IOCounters.WriteCount,
			ReadBytes:   diskIO.IOCounters.ReadBytes,
			WriteBytes:  diskIO.IOCounters.WriteBytes,
			ReadTime:    diskIO.IOCounters.ReadTime,
			WriteTime:   diskIO.IOCounters.WriteTime,
			CollectTime: collectTime,
		}

		err = monitorRepo.CreateDiskIOMonitor(diskIOMonitor)
		if err != nil {
			zap.L().Error("保存磁盘IO监控数据失败",
				zap.String("disk_name", diskIO.Device),
				zap.Error(err))
		}
	}

	zap.L().Info("磁盘IO监控数据保存成功", zap.Int("count", len(diskIOStats)))

	// 保存网卡流量监控数据
	for _, netIO := range netIOStats {
		// 跳过虚拟网卡
		if c.shouldSkipInterface(netIO.Name) {
			continue
		}

		networkMonitor := &model.NetworkMonitor{
			InterfaceName: netIO.Name,
			BytesSent:     netIO.BytesSent,
			BytesRecv:     netIO.BytesRecv,
			PacketsSent:   netIO.PacketsSent,
			PacketsRecv:   netIO.PacketsRecv,
			ErrIn:         netIO.Errin,
			ErrOut:        netIO.Errout,
			DropIn:        netIO.Dropin,
			DropOut:       netIO.Dropout,
			FIFOIn:        0, // gopsutil 不提供 FIFO 队列数据
			FIFOOut:       0,
			CollectTime:   collectTime,
		}

		err = monitorRepo.CreateNetworkMonitor(networkMonitor)
		if err != nil {
			zap.L().Error("保存网卡流量监控数据失败",
				zap.String("interface_name", netIO.Name),
				zap.Error(err))
		}
	}

	zap.L().Info("网卡流量监控数据保存成功", zap.Int("count", len(netIOStats)))

	// 删除过期的监控数据
	c.deleteExpiredMonitorData(configRepo, monitorRepo)
}

// deleteExpiredMonitorData 删除过期的监控数据
func (c *Cron) deleteExpiredMonitorData(configRepo configRepo.Repository, monitorRepo monitorRepo.Repository) {
	// 从配置中获取数据保留时长（秒）
	expiredSeconds, err := getMonitorExpired(configRepo)
	if err != nil {
		zap.L().Error("获取监控数据保留时长配置失败，使用默认值604800秒", zap.Error(err))
		expiredSeconds = 604800
	}

	// 计算过期时间
	expiredTime := time.Now().Add(-time.Duration(expiredSeconds) * time.Second)

	// 删除过期数据
	err = monitorRepo.DeleteBeforeTime(expiredTime)
	if err != nil {
		zap.L().Error("删除过期监控数据失败", zap.Error(err))
		return
	}

	zap.L().Info("删除过期监控数据成功",
		zap.Int("expired_seconds", expiredSeconds),
		zap.Time("expired_time", expiredTime),
	)
}

// shouldSkipDisk 判断是否跳过磁盘设备
func (c *Cron) shouldSkipDisk(device string) bool {
	// 跳过 loop 设备
	if strings.HasPrefix(device, "loop") {
		return true
	}
	// 跳过 zram 设备
	if strings.HasPrefix(device, "zram") {
		return true
	}
	// 跳过 dm- 设备（逻辑卷）
	if strings.HasPrefix(device, "dm-") {
		return true
	}
	return false
}

// shouldSkipInterface 判断是否跳过网络接口
func (c *Cron) shouldSkipInterface(name string) bool {
	virtualPrefixes := []string{
		"docker", "k8s", "kube", "flannel", "cni", "calico",
		"veth", "virbr", "tun", "tap", "vif", "vni",
		"br-", "ovs", "vxlan", "geneve", "gre",
		"ip_vti", "ip6tnl", "sit", "ip6gre", "lo",
	}

	nameLower := strings.ToLower(name)
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(nameLower, prefix) {
			return true
		}
	}

	return false
}

// getMonitorInterval 从配置中获取监控间隔时间
func getMonitorInterval(repo configRepo.Repository) (int, error) {
	value, err := repo.GetByKey("monitor_interval")
	if err != nil {
		return 300, err
	}
	return strconv.Atoi(value)
}

// getMonitorExpired 从配置中获取监控数据保留时长
func getMonitorExpired(repo configRepo.Repository) (int, error) {
	value, err := repo.GetByKey("monitor_expired")
	if err != nil {
		return 604800, err
	}
	return strconv.Atoi(value)
}
