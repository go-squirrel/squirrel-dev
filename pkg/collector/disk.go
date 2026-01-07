package collector

import (
	"regexp"

	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
)

type DiskInfo struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

func GetDiskInfo() (diskInfos []DiskInfo) {
	fsFilterRegex := regexp.MustCompile(`^(fuse.lxcfs|loop|nsfs|tmpfs|autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$`)
	pathFilterRegex := regexp.MustCompile(`^/(var/lib/kubelet|snap|dev|proc|sys|var/lib/docker/.+)($|/)`)
	partitions, err := disk.Partitions(true)
	if err != nil {
		zap.S().Error(err)
		return
	}

	var filteredPartitions []disk.PartitionStat
	for _, partition := range partitions {
		if !fsFilterRegex.MatchString(partition.Fstype) && !pathFilterRegex.MatchString(partition.Mountpoint) {
			filteredPartitions = append(filteredPartitions, partition)
		}
	}

	diskInfos = []DiskInfo{}
	for _, partition := range filteredPartitions {

		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // 如果无法获取某个分区的信息，则跳过
		}
		diskInfos = append(diskInfos, DiskInfo{
			Path:        partition.Mountpoint,
			Total:       usageStat.Total,
			Used:        usageStat.Used,
			Free:        usageStat.Free,
			UsedPercent: usageStat.UsedPercent,
		})
	}

	return diskInfos
}

type IoStats struct {
	ReadCount        uint64 `json:"read_count"`
	MergedReadCount  uint64 `json:"merged_read_count"`
	WriteCount       uint64 `json:"write_count"`
	MergedWriteCount uint64 `json:"merged_write_count"`
	ReadBytes        uint64 `json:"read_bytes"`
	WriteBytes       uint64 `json:"write_bytes"`
	ReadTime         uint64 `json:"read_time"`
	WriteTime        uint64 `json:"write_time"`
	IopsInProgress   uint64 `json:"iopsIn_progress"`
	IoTime           uint64 `json:"io_time"`
	WeightedIO       uint64 `json:"weighted_io"`
	Name             string `json:"name"`
	SerialNumber     string `json:"serial_number"`
	Label            string `json:"label"`
}

func GetIoStats() (stats []IoStats) {
	stats = []IoStats{}
	mapStat, _ := disk.IOCounters()
	for _, stat := range mapStat {
		stats = append(stats, IoStats(stat))
	}
	return stats
}
