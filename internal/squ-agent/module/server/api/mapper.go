package api

import (
	"squirrel-dev/internal/squ-agent/module/server/api/res"
	"squirrel-dev/internal/squ-agent/module/server/domain"
)

func toResponse(value *domain.HostInfo) *res.ServerInfo {
	if value == nil {
		return nil
	}

	addresses := make([]res.NetAddr, 0, len(value.IPAddresses))
	for _, address := range value.IPAddresses {
		addresses = append(addresses, res.NetAddr{
			Name: address.Name,
			IPv4: address.IPv4,
			IPv6: address.IPv6,
		})
	}

	return &res.ServerInfo{
		Hostname:        value.Hostname,
		OS:              value.OS,
		Platform:        value.Platform,
		PlatformVersion: value.PlatformVersion,
		KernelVersion:   value.KernelVersion,
		Architecture:    value.Architecture,
		Uptime:          value.Uptime,
		UptimeStr:       value.UptimeStr,
		IPAddresses:     addresses,
	}
}
