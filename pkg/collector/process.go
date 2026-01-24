package collector

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessStats 进程统计信息
type ProcessStats struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryMB      float64 `json:"memoryMB"`
	MemoryPercent float32 `json:"memoryPercent"`
	Status        string  `json:"status"`
	CreateTime    int64   `json:"createTime"`
	Cmdline       string  `json:"cmdline,omitempty"`
}

type Process struct {
	BaseCollector
}

func NewProcessCollector() *Process {
	return &Process{
		BaseCollector: BaseCollector{name: "process"},
	}
}

func (p *Process) Collect() (any, error) {
	return p.CollectAllProcesses()
}

// CollectTopCPU 收集CPU使用率前N的进程
func (p *Process) CollectTopCPU(limit int) ([]ProcessStats, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	var stats []ProcessStats
	for _, proc := range procs {
		stat, err := p.getProcessStats(proc)
		if err != nil {
			continue // 跳过无法获取信息的进程
		}
		stats = append(stats, *stat)
	}

	// 按CPU使用率降序排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].CPUPercent > stats[j].CPUPercent
	})

	// 返回前N个
	if limit > 0 && len(stats) > limit {
		stats = stats[:limit]
	}

	return stats, nil
}

// CollectTopMemory 收集内存使用量前N的进程
func (p *Process) CollectTopMemory(limit int) ([]ProcessStats, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	var stats []ProcessStats
	for _, proc := range procs {
		stat, err := p.getProcessStats(proc)
		if err != nil {
			continue // 跳过无法获取信息的进程
		}
		stats = append(stats, *stat)
	}

	// 按内存使用量降序排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].MemoryMB > stats[j].MemoryMB
	})

	// 返回前N个
	if limit > 0 && len(stats) > limit {
		stats = stats[:limit]
	}

	return stats, nil
}

// CollectAllProcesses 收集所有进程信息
func (p *Process) CollectAllProcesses() ([]ProcessStats, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	var stats []ProcessStats
	for _, proc := range procs {
		stat, err := p.getProcessStats(proc)
		if err != nil {
			continue
		}
		stats = append(stats, *stat)
	}

	return stats, nil
}

// getProcessStats 获取单个进程的统计信息
func (p *Process) getProcessStats(proc *process.Process) (*ProcessStats, error) {
	pid := proc.Pid

	name, err := proc.Name()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return nil, err
	}

	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return nil, err
	}

	statuses, err := proc.Status()
	if err != nil {
		statuses = []string{"unknown"}
	}

	status := "unknown"
	if len(statuses) > 0 {
		status = statuses[0]
	}

	createTime, err := proc.CreateTime()
	if err != nil {
		createTime = 0
	}

	// 计算内存百分比（可选，需要总内存信息）
	var memPercent float32 = 0

	stat := &ProcessStats{
		PID:           pid,
		Name:          name,
		CPUPercent:    cpuPercent,
		MemoryMB:      float64(memInfo.RSS) / 1024 / 1024,
		MemoryPercent: memPercent,
		Status:        status,
		CreateTime:    createTime,
	}

	return stat, nil
}
