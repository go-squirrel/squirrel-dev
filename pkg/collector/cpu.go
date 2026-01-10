package collector

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPU struct {
	BaseCollector
}

func NewCPUCollector() *CPU {
	return &CPU{
		BaseCollector: BaseCollector{name: "cpu"},
	}
}

func (c *CPU) Collect() (any, error) {
	return c.CollectCPU()
}

func (c *CPU) CollectCPU() (*CPUInfo, error) {
	info := &CPUInfo{}

	// 获取CPU核心数
	info.Cores = runtime.NumCPU()

	// 获取CPU使用率（1秒间隔）
	percent, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("获取CPU使用率失败: %w", err)
	}

	// 计算总使用率
	var totalPercent float64
	for _, p := range percent {
		totalPercent += p
	}
	info.Usage = totalPercent / float64(len(percent))
	info.PerCoreUsage = percent

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.Model = cpuInfo[0].ModelName
		info.Frequency = cpuInfo[0].Mhz
	}

	// 获取负载（仅Linux/Mac）
	info.LoadAverage = c.getLoadAverage()

	return info, nil
}

func (c *CPU) getLoadAverage() [3]float64 {
	var load [3]float64

	// 尝试读取/proc/loadavg（Linux）
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			for i := 0; i < 3; i++ {
				if val, err := strconv.ParseFloat(fields[i], 64); err == nil {
					load[i] = val
				}
			}
		}
	}

	return load
}
