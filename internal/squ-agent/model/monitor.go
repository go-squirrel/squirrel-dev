package model

import "time"

// BaseMonitor 基础监控表，包含CPU、内存、磁盘的基本监控数据
type BaseMonitor struct {
	BaseModel
	CPUUsage    float64   `gorm:"column:cpu_usage;type:decimal(5,2)" json:"cpu_usage"`                                    // CPU使用率
	MemoryUsage float64   `gorm:"column:memory_usage;type:decimal(5,2)" json:"memory_usage"`                              // 内存使用率
	MemoryTotal uint64    `gorm:"column:memory_total;type:bigint unsigned" json:"memory_total"`                           // 内存总量(bytes)
	MemoryUsed  uint64    `gorm:"column:memory_used;type:bigint unsigned" json:"memory_used"`                             // 内存已用量(bytes)
	DiskUsage   float64   `gorm:"column:disk_usage;type:decimal(5,2)" json:"disk_usage"`                                  // 磁盘使用率
	DiskTotal   uint64    `gorm:"column:disk_total;type:bigint unsigned" json:"disk_total"`                               // 磁盘总量(bytes)
	DiskUsed    uint64    `gorm:"column:disk_used;type:bigint unsigned" json:"disk_used"`                                 // 磁盘已用量(bytes)
	CollectTime time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index" json:"collect_time"` // 数据收集时间
}

// DiskIOMonitor 磁盘IO表，记录磁盘IO数据
type DiskIOMonitor struct {
	BaseModel
	DiskName       string    `gorm:"column:disk_name;type:varchar(100);not null" json:"disk_name"`                           // 磁盘名称
	ReadCount      uint64    `gorm:"column:read_count;type:bigint unsigned" json:"read_count"`                               // 读取次数
	WriteCount     uint64    `gorm:"column:write_count;type:bigint unsigned" json:"write_count"`                             // 写入次数
	ReadBytes      uint64    `gorm:"column:read_bytes;type:bigint unsigned" json:"read_bytes"`                               // 读取字节数
	WriteBytes     uint64    `gorm:"column:write_bytes;type:bigint unsigned" json:"write_bytes"`                             // 写入字节数
	ReadTime       uint64    `gorm:"column:read_time;type:bigint unsigned" json:"read_time"`                                 // 读取耗时(ms)
	WriteTime      uint64    `gorm:"column:write_time;type:bigint unsigned" json:"write_time"`                               // 写入耗时(ms)
	IoTime         uint64    `gorm:"column:io_time;type:bigint unsigned" json:"io_time"`                                     // IO总耗时(ms)
	WeightedIoTime uint64    `gorm:"column:weighted_io_time;type:bigint unsigned" json:"weighted_io_time"`                   // 加权IO时间(ms)
	IopsInProgress uint64    `gorm:"column:iops_in_progress;type:bigint unsigned" json:"iops_in_progress"`                   // 正在处理的IOPS数
	CollectTime    time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index" json:"collect_time"` // 数据收集时间
}

// NetworkMonitor 网卡流量表，记录网络流量数据
type NetworkMonitor struct {
	BaseModel
	InterfaceName string    `gorm:"column:interface_name;type:varchar(100);not null" json:"interface_name"`                 // 网络接口名称
	BytesSent     uint64    `gorm:"column:bytes_sent;type:bigint unsigned" json:"bytes_sent"`                               // 发送字节数
	BytesRecv     uint64    `gorm:"column:bytes_recv;type:bigint unsigned" json:"bytes_recv"`                               // 接收字节数
	PacketsSent   uint64    `gorm:"column:packets_sent;type:bigint unsigned" json:"packets_sent"`                           // 发送包数
	PacketsRecv   uint64    `gorm:"column:packets_recv;type:bigint unsigned" json:"packets_recv"`                           // 接收包数
	ErrIn         uint64    `gorm:"column:err_in;type:bigint unsigned" json:"err_in"`                                       // 接收错误数
	ErrOut        uint64    `gorm:"column:err_out;type:bigint unsigned" json:"err_out"`                                     // 发送错误数
	DropIn        uint64    `gorm:"column:drop_in;type:bigint unsigned" json:"drop_in"`                                     // 接收丢弃数
	DropOut       uint64    `gorm:"column:drop_out;type:bigint unsigned" json:"drop_out"`                                   // 发送丢弃数
	FIFOIn        uint64    `gorm:"column:fifo_in;type:bigint unsigned" json:"fifo_in"`                                     // 接收FIFO队列数
	FIFOOut       uint64    `gorm:"column:fifo_out;type:bigint unsigned" json:"fifo_out"`                                   // 发送FIFO队列数
	CollectTime   time.Time `gorm:"column:collect_time;type:timestamp;default:CURRENT_TIMESTAMP;index" json:"collect_time"` // 数据收集时间
}
