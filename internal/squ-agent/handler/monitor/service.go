package monitor

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	monitorres "squirrel-dev/internal/squ-agent/handler/monitor/res"
	"squirrel-dev/internal/squ-agent/model"
	monitorRepository "squirrel-dev/internal/squ-agent/repository/monitor"
	"squirrel-dev/pkg/collector"

	"go.uber.org/zap"
)

type Monitor struct {
	Config     *config.Config
	Repository monitorRepository.Repository
	Factory    *collector.CollectorFactory
}

// GetBaseMonitorPage 获取基础监控数据分页
func (m *Monitor) GetBaseMonitorPage(page, count int) response.Response {
	monitors, total, err := m.Repository.GetBaseMonitorPage(page, count)
	if err != nil {
		zap.L().Error("Failed to get base monitor page", zap.Int("page", page), zap.Int("count", count), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 转换为响应结构
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

	result := monitorres.PageData{
		List:  responseList,
		Total: total,
		Page:  page,
		Size:  count,
	}

	return response.Success(result)
}

// GetDiskIOMonitorPage 获取磁盘IO监控数据分页
func (m *Monitor) GetDiskIOMonitorPage(page, count int) response.Response {
	monitors, total, err := m.Repository.GetDiskIOMonitorPage(page, count)
	if err != nil {
		zap.L().Error("Failed to get disk IO monitor page", zap.Int("page", page), zap.Int("count", count), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 转换为响应结构
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

	result := monitorres.PageData{
		List:  responseList,
		Total: total,
		Page:  page,
		Size:  count,
	}

	return response.Success(result)
}

// GetNetworkMonitorPage 获取网卡流量监控数据分页
func (m *Monitor) GetNetworkMonitorPage(page, count int) response.Response {
	monitors, total, err := m.Repository.GetNetworkMonitorPage(page, count)
	if err != nil {
		zap.L().Error("Failed to get network monitor page", zap.Int("page", page), zap.Int("count", count), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 转换为响应结构
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

	result := monitorres.PageData{
		List:  responseList,
		Total: total,
		Page:  page,
		Size:  count,
	}

	return response.Success(result)
}
