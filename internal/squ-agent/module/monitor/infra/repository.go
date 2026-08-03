package infra

import (
	"context"
	"time"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-agent/module/monitor/domain"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateBaseMonitor(ctx context.Context, value *domain.BaseMonitor) error {
	model := baseToModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) CreateDiskIOMonitor(ctx context.Context, value *domain.DiskIOMonitor) error {
	model := diskIOToModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) CreateNetworkMonitor(ctx context.Context, value *domain.NetworkMonitor) error {
	model := networkToModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) CreateDiskUsageMonitor(ctx context.Context, value *domain.DiskUsageMonitor) error {
	model := diskUsageToModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) DeleteBeforeTime(ctx context.Context, before time.Time) error {
	db := r.db.WithContext(ctx)
	for _, model := range []any{&baseMonitorModel{}, &diskIOMonitorModel{}, &networkMonitorModel{}, &diskUsageMonitorModel{}} {
		if err := db.Where("collect_time < ?", before).Delete(model).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) BaseByTimeRange(ctx context.Context, since time.Time) ([]domain.BaseMonitor, error) {
	var models []baseMonitorModel
	if err := r.query(ctx, since).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.BaseMonitor
	for _, model := range models {
		result = append(result, baseToDomain(model))
	}
	return result, nil
}

func (r *Repository) DiskIOByTimeRange(ctx context.Context, since time.Time) ([]domain.DiskIOMonitor, error) {
	var models []diskIOMonitorModel
	if err := r.query(ctx, since).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.DiskIOMonitor
	for _, model := range models {
		result = append(result, diskIOToDomain(model))
	}
	return result, nil
}

func (r *Repository) NetworkByTimeRange(ctx context.Context, since time.Time) ([]domain.NetworkMonitor, error) {
	var models []networkMonitorModel
	if err := r.query(ctx, since).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.NetworkMonitor
	for _, model := range models {
		result = append(result, networkToDomain(model))
	}
	return result, nil
}

func (r *Repository) DiskUsageByTimeRange(ctx context.Context, since time.Time) ([]domain.DiskUsageMonitor, error) {
	var models []diskUsageMonitorModel
	if err := r.query(ctx, since).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.DiskUsageMonitor
	for _, model := range models {
		result = append(result, diskUsageToDomain(model))
	}
	return result, nil
}

func (r *Repository) query(ctx context.Context, since time.Time) *gorm.DB {
	return r.db.WithContext(ctx).Where("collect_time >= ?", since).Order("collect_time ASC")
}

func baseToModel(v domain.BaseMonitor) baseMonitorModel {
	return baseMonitorModel{baseModel: baseModel{ID: v.ID}, CPUUsage: v.CPUUsage, MemoryUsage: v.MemoryUsage, MemoryTotal: v.MemoryTotal, MemoryUsed: v.MemoryUsed, DiskUsage: v.DiskUsage, DiskTotal: v.DiskTotal, DiskUsed: v.DiskUsed, CollectTime: v.CollectTime}
}
func baseToDomain(v baseMonitorModel) domain.BaseMonitor {
	return domain.BaseMonitor{ID: v.ID, CPUUsage: v.CPUUsage, MemoryUsage: v.MemoryUsage, MemoryTotal: v.MemoryTotal, MemoryUsed: v.MemoryUsed, DiskUsage: v.DiskUsage, DiskTotal: v.DiskTotal, DiskUsed: v.DiskUsed, CollectTime: v.CollectTime}
}
func diskIOToModel(v domain.DiskIOMonitor) diskIOMonitorModel {
	return diskIOMonitorModel{baseModel: baseModel{ID: v.ID}, DiskName: v.DiskName, ReadCount: v.ReadCount, WriteCount: v.WriteCount, ReadBytes: v.ReadBytes, WriteBytes: v.WriteBytes, ReadTime: v.ReadTime, WriteTime: v.WriteTime, IoTime: v.IoTime, WeightedIoTime: v.WeightedIoTime, IopsInProgress: v.IopsInProgress, CollectTime: v.CollectTime}
}
func diskIOToDomain(v diskIOMonitorModel) domain.DiskIOMonitor {
	return domain.DiskIOMonitor{ID: v.ID, DiskName: v.DiskName, ReadCount: v.ReadCount, WriteCount: v.WriteCount, ReadBytes: v.ReadBytes, WriteBytes: v.WriteBytes, ReadTime: v.ReadTime, WriteTime: v.WriteTime, IoTime: v.IoTime, WeightedIoTime: v.WeightedIoTime, IopsInProgress: v.IopsInProgress, CollectTime: v.CollectTime}
}
func networkToModel(v domain.NetworkMonitor) networkMonitorModel {
	return networkMonitorModel{baseModel: baseModel{ID: v.ID}, InterfaceName: v.InterfaceName, BytesSent: v.BytesSent, BytesRecv: v.BytesRecv, PacketsSent: v.PacketsSent, PacketsRecv: v.PacketsRecv, ErrIn: v.ErrIn, ErrOut: v.ErrOut, DropIn: v.DropIn, DropOut: v.DropOut, FIFOIn: v.FIFOIn, FIFOOut: v.FIFOOut, CollectTime: v.CollectTime}
}
func networkToDomain(v networkMonitorModel) domain.NetworkMonitor {
	return domain.NetworkMonitor{ID: v.ID, InterfaceName: v.InterfaceName, BytesSent: v.BytesSent, BytesRecv: v.BytesRecv, PacketsSent: v.PacketsSent, PacketsRecv: v.PacketsRecv, ErrIn: v.ErrIn, ErrOut: v.ErrOut, DropIn: v.DropIn, DropOut: v.DropOut, FIFOIn: v.FIFOIn, FIFOOut: v.FIFOOut, CollectTime: v.CollectTime}
}
func diskUsageToModel(v domain.DiskUsageMonitor) diskUsageMonitorModel {
	return diskUsageMonitorModel{baseModel: baseModel{ID: v.ID}, DeviceName: v.DeviceName, MountPoint: v.MountPoint, FsType: v.FsType, Total: v.Total, Used: v.Used, Free: v.Free, Usage: v.Usage, InodesTotal: v.InodesTotal, InodesUsed: v.InodesUsed, InodesFree: v.InodesFree, CollectTime: v.CollectTime}
}
func diskUsageToDomain(v diskUsageMonitorModel) domain.DiskUsageMonitor {
	return domain.DiskUsageMonitor{ID: v.ID, DeviceName: v.DeviceName, MountPoint: v.MountPoint, FsType: v.FsType, Total: v.Total, Used: v.Used, Free: v.Free, Usage: v.Usage, InodesTotal: v.InodesTotal, InodesUsed: v.InodesUsed, InodesFree: v.InodesFree, CollectTime: v.CollectTime}
}
