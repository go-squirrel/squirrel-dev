package infra

import (
	"time"

	"gorm.io/gorm"
)

type baseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type baseMonitorModel struct {
	baseModel
	CPUUsage    float64   `gorm:"column:cpu_usage;type:decimal(5,2)"`
	MemoryUsage float64   `gorm:"column:memory_usage;type:decimal(5,2)"`
	MemoryTotal uint64    `gorm:"column:memory_total;type:bigint unsigned"`
	MemoryUsed  uint64    `gorm:"column:memory_used;type:bigint unsigned"`
	DiskUsage   float64   `gorm:"column:disk_usage;type:decimal(5,2)"`
	DiskTotal   uint64    `gorm:"column:disk_total;type:bigint unsigned"`
	DiskUsed    uint64    `gorm:"column:disk_used;type:bigint unsigned"`
	CollectTime time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (baseMonitorModel) TableName() string { return "base_monitors" }

type diskIOMonitorModel struct {
	baseModel
	DiskName       string    `gorm:"column:disk_name;type:varchar(100);not null"`
	ReadCount      uint64    `gorm:"column:read_count;type:bigint unsigned"`
	WriteCount     uint64    `gorm:"column:write_count;type:bigint unsigned"`
	ReadBytes      uint64    `gorm:"column:read_bytes;type:bigint unsigned"`
	WriteBytes     uint64    `gorm:"column:write_bytes;type:bigint unsigned"`
	ReadTime       uint64    `gorm:"column:read_time;type:bigint unsigned"`
	WriteTime      uint64    `gorm:"column:write_time;type:bigint unsigned"`
	IoTime         uint64    `gorm:"column:io_time;type:bigint unsigned"`
	WeightedIoTime uint64    `gorm:"column:weighted_io_time;type:bigint unsigned"`
	IopsInProgress uint64    `gorm:"column:iops_in_progress;type:bigint unsigned"`
	CollectTime    time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (diskIOMonitorModel) TableName() string { return "disk_io_monitors" }

type diskUsageMonitorModel struct {
	baseModel
	DeviceName  string    `gorm:"column:device_name;type:varchar(100);not null"`
	MountPoint  string    `gorm:"column:mount_point;type:varchar(255)"`
	FsType      string    `gorm:"column:fs_type;type:varchar(50)"`
	Total       uint64    `gorm:"column:total;type:bigint unsigned"`
	Used        uint64    `gorm:"column:used;type:bigint unsigned"`
	Free        uint64    `gorm:"column:free;type:bigint unsigned"`
	Usage       float64   `gorm:"column:usage;type:decimal(5,2)"`
	InodesTotal uint64    `gorm:"column:inodes_total;type:bigint unsigned"`
	InodesUsed  uint64    `gorm:"column:inodes_used;type:bigint unsigned"`
	InodesFree  uint64    `gorm:"column:inodes_free;type:bigint unsigned"`
	CollectTime time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (diskUsageMonitorModel) TableName() string { return "disk_usage_monitors" }

type networkMonitorModel struct {
	baseModel
	InterfaceName string    `gorm:"column:interface_name;type:varchar(100);not null"`
	BytesSent     uint64    `gorm:"column:bytes_sent;type:bigint unsigned"`
	BytesRecv     uint64    `gorm:"column:bytes_recv;type:bigint unsigned"`
	PacketsSent   uint64    `gorm:"column:packets_sent;type:bigint unsigned"`
	PacketsRecv   uint64    `gorm:"column:packets_recv;type:bigint unsigned"`
	ErrIn         uint64    `gorm:"column:err_in;type:bigint unsigned"`
	ErrOut        uint64    `gorm:"column:err_out;type:bigint unsigned"`
	DropIn        uint64    `gorm:"column:drop_in;type:bigint unsigned"`
	DropOut       uint64    `gorm:"column:drop_out;type:bigint unsigned"`
	FIFOIn        uint64    `gorm:"column:fifo_in;type:bigint unsigned"`
	FIFOOut       uint64    `gorm:"column:fifo_out;type:bigint unsigned"`
	CollectTime   time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (networkMonitorModel) TableName() string { return "network_monitors" }
