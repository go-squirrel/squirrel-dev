package monitor

import (
	"fmt"
	"squirrel-dev/internal/pkg/cache"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	monitorres "squirrel-dev/internal/squ-agent/handler/monitor/res"
	"squirrel-dev/internal/squ-agent/model"
	monitorRepository "squirrel-dev/internal/squ-agent/repository/monitor"
	"squirrel-dev/pkg/collector"
	"time"

	"go.uber.org/zap"
)

type Monitor struct {
	Config     *config.Config
	Cache      cache.Cache
	Repository monitorRepository.Repository
	Factory    *collector.CollectorFactory
}

func New(config *config.Config, cache cache.Cache, repo monitorRepository.Repository, factory *collector.CollectorFactory) *Monitor {
	return &Monitor{
		Config:     config,
		Cache:      cache,
		Repository: repo,
		Factory:    factory,
	}
}

// parseTimeRange 解析时间范围字符串
func parseTimeRange(rangeStr string) (time.Time, error) {
	now := time.Now()
	switch rangeStr {
	case "1h":
		return now.Add(-1 * time.Hour), nil
	case "6h":
		return now.Add(-6 * time.Hour), nil
	case "24h":
		return now.Add(-24 * time.Hour), nil
	case "7d":
		return now.Add(-7 * 24 * time.Hour), nil
	default:
		return time.Time{}, fmt.Errorf("invalid time range: %s", rangeStr)
	}
}

// GetBaseMonitorByRange 按时间范围获取基础监控数据
func (m *Monitor) GetBaseMonitorByRange(timeRange string) response.Response {
	since, err := parseTimeRange(timeRange)
	if err != nil {
		zap.L().Warn("Invalid time range", zap.String("range", timeRange), zap.Error(err))
		return response.Error(monitorres.ErrCodeParameter)
	}

	monitors, err := m.Repository.GetBaseMonitorByTimeRange(since)
	if err != nil {
		zap.L().Error("Failed to get base monitor by range",
			zap.String("range", timeRange), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	var responseList []monitorres.BaseMonitorResponse
	for _, item := range monitors {
		responseList = append(responseList, monitorres.BaseMonitorResponse{
			ID:          item.ID,
			CPUUsage:    item.CPUUsage,
			MemoryUsage: item.MemoryUsage,
			MemoryTotal: item.MemoryTotal,
			MemoryUsed:  item.MemoryUsed,
			DiskUsage:   item.DiskUsage,
			DiskTotal:   item.DiskTotal,
			DiskUsed:    item.DiskUsed,
			CollectTime: item.CollectTime,
		})
	}

	return response.Success(responseList)
}

// GetDiskIOMonitorByRange 按时间范围获取磁盘IO监控数据
func (m *Monitor) GetDiskIOMonitorByRange(timeRange string) response.Response {
	since, err := parseTimeRange(timeRange)
	if err != nil {
		zap.L().Warn("Invalid time range", zap.String("range", timeRange), zap.Error(err))
		return response.Error(monitorres.ErrCodeParameter)
	}

	monitors, err := m.Repository.GetDiskIOMonitorByTimeRange(since)
	if err != nil {
		zap.L().Error("Failed to get disk IO monitor by range",
			zap.String("range", timeRange), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	var responseList []monitorres.DiskIOMonitorResponse
	for _, item := range monitors {
		responseList = append(responseList, monitorres.DiskIOMonitorResponse{
			ID:             item.ID,
			DiskName:       item.DiskName,
			ReadCount:      item.ReadCount,
			WriteCount:     item.WriteCount,
			ReadBytes:      item.ReadBytes,
			WriteBytes:     item.WriteBytes,
			ReadTime:       item.ReadTime,
			WriteTime:      item.WriteTime,
			IoTime:         item.IoTime,
			WeightedIoTime: item.WeightedIoTime,
			IopsInProgress: item.IopsInProgress,
			CollectTime:    item.CollectTime,
		})
	}

	return response.Success(responseList)
}

// GetDiskUsageMonitorByRange 按时间范围获取磁盘使用监控数据
func (m *Monitor) GetDiskUsageMonitorByRange(timeRange string) response.Response {
	since, err := parseTimeRange(timeRange)
	if err != nil {
		zap.L().Warn("Invalid time range", zap.String("range", timeRange), zap.Error(err))
		return response.Error(monitorres.ErrCodeParameter)
	}

	monitors, err := m.Repository.GetDiskUsageMonitorByTimeRange(since)
	if err != nil {
		zap.L().Error("Failed to get disk usage monitor by range",
			zap.String("range", timeRange), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	var responseList []monitorres.DiskUsageMonitorResponse
	for _, item := range monitors {
		responseList = append(responseList, monitorres.DiskUsageMonitorResponse{
			ID:          item.ID,
			DeviceName:  item.DeviceName,
			MountPoint:  item.MountPoint,
			FsType:      item.FsType,
			Total:       item.Total,
			Used:        item.Used,
			Free:        item.Free,
			Usage:       item.Usage,
			InodesTotal: item.InodesTotal,
			InodesUsed:  item.InodesUsed,
			InodesFree:  item.InodesFree,
			CollectTime: item.CollectTime,
		})
	}

	return response.Success(responseList)
}

// GetNetworkMonitorByRange 按时间范围获取网络监控数据
func (m *Monitor) GetNetworkMonitorByRange(timeRange string) response.Response {
	since, err := parseTimeRange(timeRange)
	if err != nil {
		zap.L().Warn("Invalid time range", zap.String("range", timeRange), zap.Error(err))
		return response.Error(monitorres.ErrCodeParameter)
	}

	monitors, err := m.Repository.GetNetworkMonitorByTimeRange(since)
	if err != nil {
		zap.L().Error("Failed to get network monitor by range",
			zap.String("range", timeRange), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	var responseList []monitorres.NetworkMonitorResponse
	for _, item := range monitors {
		responseList = append(responseList, monitorres.NetworkMonitorResponse{
			ID:            item.ID,
			InterfaceName: item.InterfaceName,
			BytesSent:     item.BytesSent,
			BytesRecv:     item.BytesRecv,
			PacketsSent:   item.PacketsSent,
			PacketsRecv:   item.PacketsRecv,
			ErrIn:         item.ErrIn,
			ErrOut:        item.ErrOut,
			DropIn:        item.DropIn,
			DropOut:       item.DropOut,
			FIFOIn:        item.FIFOIn,
			FIFOOut:       item.FIFOOut,
			CollectTime:   item.CollectTime,
		})
	}

	return response.Success(responseList)
}
