package server

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/server/res"
	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/collector"

	"go.uber.org/zap"
)

type Server struct {
	Config  *config.Config
	Factory *collector.CollectorFactory
}

func (s *Server) GetInfo() response.Response {
	if s.Factory == nil {
		zap.L().Error("Factory is nil")
		return response.Error(model.ReturnErrCode(nil))
	}

	hostCollector := s.Factory.GetHostCollector()
	if hostCollector == nil {
		zap.L().Error("Host collector is nil")
		return response.Error(model.ReturnErrCode(nil))
	}

	hostInfo, err := hostCollector.CollectHostInfo()
	if err != nil {
		zap.L().Error("Failed to collect host info", zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 转换响应格式
	var netAddrs []res.NetAddr
	for _, addr := range hostInfo.IPAddresses {
		netAddrs = append(netAddrs, res.NetAddr{
			Name: addr.Name,
			IPv4: addr.IPv4,
			IPv6: addr.IPv6,
		})
	}

	info := res.ServerInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		Architecture:    hostInfo.Architecture,
		Uptime:          hostInfo.Uptime,
		UptimeStr:       hostInfo.UptimeStr,
		IPAddresses:     netAddrs,
	}

	return response.Success(info)
}
