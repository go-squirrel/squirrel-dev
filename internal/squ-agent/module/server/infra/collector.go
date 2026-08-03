package infra

import (
	"context"

	"squirrel-dev/internal/squ-agent/module/server/domain"
	"squirrel-dev/pkg/collector"
)

type hostInfoCollector struct {
	source collector.HostCollector
}

func NewHostInfoCollector(source collector.HostCollector) domain.HostInfoCollector {
	return &hostInfoCollector{source: source}
}

func (c *hostInfoCollector) CollectHostInfo(_ context.Context) (*domain.HostInfo, error) {
	value, err := c.source.CollectHostInfo()
	if err != nil {
		return nil, err
	}

	addresses := make([]domain.NetAddr, 0, len(value.IPAddresses))
	for _, address := range value.IPAddresses {
		addresses = append(addresses, domain.NetAddr{
			Name: address.Name,
			IPv4: address.IPv4,
			IPv6: address.IPv6,
		})
	}

	return &domain.HostInfo{
		Hostname:        value.Hostname,
		OS:              value.OS,
		Platform:        value.Platform,
		PlatformVersion: value.PlatformVersion,
		KernelVersion:   value.KernelVersion,
		Architecture:    value.Architecture,
		Uptime:          value.Uptime,
		UptimeStr:       value.UptimeStr,
		IPAddresses:     addresses,
	}, nil
}
