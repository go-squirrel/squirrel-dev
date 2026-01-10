package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type Memory struct {
	BaseCollector
}

func NewMemoryCollector() *Memory {
	return &Memory{
		BaseCollector: BaseCollector{name: "memory"},
	}
}

func (m *Memory) Collect() (any, error) {
	return m.CollectMemory()
}

func (m *Memory) CollectMemory() (*MemInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("获取内存信息失败: %w", err)
	}

	swapStat, err := mem.SwapMemory()
	if err != nil {
		// 交换分区可能不存在，只记录错误但不返回失败
		fmt.Printf("获取交换分区信息失败: %v\n", err)
	}

	info := &MemInfo{
		Total:     vmStat.Total,
		Available: vmStat.Available,
		Used:      vmStat.Used,
		Usage:     vmStat.UsedPercent,
	}

	if swapStat != nil {
		info.SwapTotal = swapStat.Total
		info.SwapUsed = swapStat.Used
	}

	return info, nil
}
