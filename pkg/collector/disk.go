package collector

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

type Disk struct {
	BaseCollector
}

func NewDiskCollector() *Disk {
	return &Disk{
		BaseCollector: BaseCollector{name: "disk"},
	}
}

func (d *Disk) Collect() (any, error) {
	return d.CollectDisk()
}

func (d *Disk) CollectDisk() (*DiskInfo, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("获取磁盘分区失败: %w", err)
	}

	diskInfo := &DiskInfo{}
	var total, used, free uint64

	for _, partition := range partitions {
		// 跳过特殊文件系统
		if d.shouldSkip(partition.Fstype) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // 跳过无法访问的分区
		}

		partitionInfo := Partition{
			Device:      partition.Device,
			MountPoint:  partition.Mountpoint,
			FSType:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Available:   usage.Free,
			Usage:       usage.UsedPercent,
			InodesTotal: usage.InodesTotal,
			InodesUsed:  usage.InodesUsed,
			InodesFree:  usage.InodesFree,
		}

		diskInfo.Partitions = append(diskInfo.Partitions, partitionInfo)

		// 累计总量
		total += usage.Total
		used += usage.Used
		free += usage.Free
	}

	// 计算总体使用率
	diskInfo.Total = total
	diskInfo.Used = used
	diskInfo.Available = free
	if total > 0 {
		diskInfo.Usage = (float64(used) / float64(total)) * 100
	}

	return diskInfo, nil
}

func (d *Disk) shouldSkip(fsType string) bool {
	// 跳过虚拟文件系统
	skipTypes := []string{
		"fuse.lxcfs", "loop", "nsfs", "tmpfs", "squashfs", "fuse",
		"autofs", "binfmt_misc", "cgroup", "cgroup2", "efivarfs", "vfat",
		"configfs", "debugfs", "devpts", "devtmpfs", "fusectl",
		"hugetlbfs", "mqueue", "overlay", "proc", "procfs", "bpf",
		"pstore", "rpc_pipefs", "securityfs", "sysfs", "tracefs",
		"var/lib/kubelet", "snap", "dev", "proc", "sys", "var/lib/docker/",
	}

	for _, t := range skipTypes {
		if strings.Contains(fsType, t) {
			return true
		}
	}
	return false
}
