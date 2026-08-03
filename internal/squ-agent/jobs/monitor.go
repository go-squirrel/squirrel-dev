package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	monitorDomain "squirrel-dev/internal/squ-agent/module/monitor/domain"
	monitorInfra "squirrel-dev/internal/squ-agent/module/monitor/infra"
	"squirrel-dev/pkg/collector"
)

func (j *Jobs) registerMonitorCollection() error {
	interval, err := j.configInt(context.Background(), "monitor_interval")
	if err != nil || interval == 0 {
		interval = 300
	}
	// Preserve the old seconds-to-whole-minutes conversion, including the
	// invalid */0 expression for configured values below 60 seconds.
	expression := fmt.Sprintf("0 */%d * * * *", interval/60)
	_, err = j.cron.AddFunc(expression, j.collectAndSaveMonitorData)
	return err
}

func (j *Jobs) collectAndSaveMonitorData() {
	collectedAt := time.Now()
	j.collectBase(collectedAt)
	j.collectDiskUsage(collectedAt)
	j.collectDiskIO(collectedAt)
	j.collectNetwork(collectedAt)
	j.deleteExpiredMonitorData()
}

func (j *Jobs) collectBase(collectedAt time.Time) {
	cpu, err := collector.NewCPUCollector().CollectCPU()
	if err != nil {
		return
	}
	memory, err := collector.NewMemoryCollector().CollectMemory()
	if err != nil {
		return
	}
	disk, err := collector.NewDiskCollector().CollectDisk()
	if err != nil {
		return
	}
	_ = j.monitors.CreateBaseMonitor(context.Background(), &monitorDomain.BaseMonitor{
		CPUUsage: cpu.Usage, MemoryUsage: memory.Usage, MemoryTotal: memory.Total,
		MemoryUsed: memory.Used, DiskUsage: disk.Usage, DiskTotal: disk.Total,
		DiskUsed: disk.Used, CollectTime: collectedAt,
	})
}

func (j *Jobs) collectDiskUsage(collectedAt time.Time) {
	disk, err := collector.NewDiskCollector().CollectDisk()
	if err != nil {
		return
	}
	for _, partition := range disk.Partitions {
		_ = j.monitors.CreateDiskUsageMonitor(context.Background(), &monitorDomain.DiskUsageMonitor{
			DeviceName: partition.Device, MountPoint: partition.MountPoint, FsType: partition.FSType,
			Total: partition.Total, Used: partition.Used, Free: partition.Available, Usage: partition.Usage,
			InodesTotal: partition.InodesTotal, InodesUsed: partition.InodesUsed,
			InodesFree: partition.InodesFree, CollectTime: collectedAt,
		})
	}
}

func (j *Jobs) collectDiskIO(collectedAt time.Time) {
	values, err := collector.NewIOCollector().CollectAllDiskIO()
	if err != nil {
		return
	}
	for _, value := range values {
		if skipDisk(value.Device) {
			continue
		}
		_ = j.monitors.CreateDiskIOMonitor(context.Background(), &monitorDomain.DiskIOMonitor{
			DiskName: value.Device, ReadCount: value.IOCounters.ReadCount,
			WriteCount: value.IOCounters.WriteCount, ReadBytes: value.IOCounters.ReadBytes,
			WriteBytes: value.IOCounters.WriteBytes, ReadTime: value.IOCounters.ReadTime,
			WriteTime: value.IOCounters.WriteTime, CollectTime: collectedAt,
		})
	}
}

func (j *Jobs) collectNetwork(collectedAt time.Time) {
	values, err := collector.NewIOCollector().CollectAllNetIO()
	if err != nil {
		return
	}
	for _, value := range values {
		if skipInterface(value.Name) {
			continue
		}
		_ = j.monitors.CreateNetworkMonitor(context.Background(), &monitorDomain.NetworkMonitor{
			InterfaceName: value.Name, BytesSent: value.BytesSent, BytesRecv: value.BytesRecv,
			PacketsSent: value.PacketsSent, PacketsRecv: value.PacketsRecv,
			ErrIn: value.Errin, ErrOut: value.Errout, DropIn: value.Dropin, DropOut: value.Dropout,
			FIFOIn: 0, FIFOOut: 0, CollectTime: collectedAt,
		})
	}
}

func (j *Jobs) deleteExpiredMonitorData() {
	expired, err := j.configInt(context.Background(), "monitor_expired")
	if err != nil || expired == 0 {
		expired = 604800
	}
	_ = j.monitors.DeleteBeforeTime(
		context.Background(),
		time.Now().Add(-time.Duration(expired)*time.Second),
	)
}

func (j *Jobs) refreshMonitorStatsCache() {
	if j.cache == nil {
		return
	}
	ctx := context.Background()
	factory := collector.NewCollectorFactory()
	factory.Register(collector.NewCPUCollector())
	factory.Register(collector.NewMemoryCollector())
	factory.Register(collector.NewDiskCollector())
	factory.Register(collector.NewIOCollector())
	factory.Register(collector.NewProcessCollector())

	host, err := factory.CollectAll()
	if err != nil {
		return
	}
	if host.CPU.Usage > 0 || len(host.CPU.PerCoreUsage) > 0 {
		_ = j.cache.Set(ctx, "monitor:stats:cpu", host.CPU, 5*time.Second)
	}
	_ = j.cache.Set(ctx, "monitor:stats:memory", host.Memory, 10*time.Second)
	_ = j.cache.Set(ctx, "monitor:stats:disk", host.Disk, 60*time.Second)

	var topCPU, topMemory []collector.ProcessStats
	if process := factory.GetProcessCollector(); process != nil {
		topCPU, _ = process.CollectTopCPU(5)
		topMemory, _ = process.CollectTopMemory(5)
	}
	if len(topCPU) > 0 || len(topMemory) > 0 {
		_ = j.cache.Set(ctx, "monitor:stats:process", processCacheData{
			TopCPU: topCPU, TopMemory: topMemory,
		}, 10*time.Second)
	}
	stats := monitorInfra.BuildStats(host, topCPU, topMemory)
	_ = j.cache.Set(ctx, monitorDomain.StatsCacheKey, stats, monitorDomain.StatsCacheTTL)
}

type processCacheData struct {
	TopCPU    []collector.ProcessStats
	TopMemory []collector.ProcessStats
}

func skipDisk(device string) bool {
	return strings.HasPrefix(device, "loop") ||
		strings.HasPrefix(device, "zram") ||
		strings.HasPrefix(device, "dm-")
}

func skipInterface(name string) bool {
	prefixes := []string{
		"docker", "k8s", "kube", "flannel", "cni", "calico",
		"veth", "virbr", "tun", "tap", "vif", "vni",
		"br-", "ovs", "vxlan", "geneve", "gre",
		"ip_vti", "ip6tnl", "sit", "ip6gre", "lo",
	}
	name = strings.ToLower(name)
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
