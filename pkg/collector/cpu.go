package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type CpuStats struct {
	TotalUsage   float64   `json:"total_usage"`
	PerCoreUsage []float64 `json:"per_core_usage"`
	Counts       int       `json:"counts"`
}

func GetCpuLinuxStats() (info CpuStats) {
	cpuCount := getCPUCount()
	totalUsage, perCoreUsage := getCPUUsage(cpuCount)
	info.Counts = cpuCount
	info.TotalUsage = totalUsage
	info.PerCoreUsage = perCoreUsage
	return info
}

func getCPUCount() int {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		zap.S().Error("open err:", err)
		return 0
	}
	defer file.Close()

	cpuCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "processor") {
			cpuCount++
		}
	}
	return cpuCount
}

func getCPUUsage(cpuCount int) (float64, []float64) {
	initialStats := readCPUStats(cpuCount)
	time.Sleep(1 * time.Second) // 等待1秒以获取CPU使用情况的样本
	finalStats := readCPUStats(cpuCount)

	totalUsage := calculateTotalUsage(initialStats, finalStats)
	perCoreUsage := calculatePerCoreUsage(initialStats, finalStats)

	return totalUsage, perCoreUsage
}

func readCPUStats(cpuCount int) [][]int64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		zap.S().Error("open err:", err)
		return nil
	}
	defer file.Close()

	var stats [][]int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)[1:]
			var stat []int64
			for _, field := range fields {
				val, _ := strconv.ParseInt(field, 10, 64)
				stat = append(stat, val)
			}
			stats = append(stats, stat)
			if len(stats) == cpuCount+1 { // 包括总体CPU使用情况和每个核心的使用情况
				break
			}
		}
	}
	return stats
}

func calculateTotalUsage(initial, final [][]int64) float64 {
	initialTotal := sum(initial[0])
	finalTotal := sum(final[0])
	initialIdle := initial[0][3]
	finalIdle := final[0][3]

	totalDelta := finalTotal - initialTotal
	idleDelta := finalIdle - initialIdle

	usage := (1.0 - float64(idleDelta)/float64(totalDelta)) * 100
	return usage
}

func calculatePerCoreUsage(initial, final [][]int64) []float64 {
	var perCoreUsage []float64
	for i := 1; i < len(initial); i++ {
		initialTotal := sum(initial[i])
		finalTotal := sum(final[i])
		initialIdle := initial[i][3]
		finalIdle := final[i][3]

		totalDelta := finalTotal - initialTotal
		idleDelta := finalIdle - initialIdle

		usage := (1.0 - float64(idleDelta)/float64(totalDelta)) * 100
		perCoreUsage = append(perCoreUsage, usage)
	}
	return perCoreUsage
}

func sum(slice []int64) int64 {
	total := int64(0)
	for _, val := range slice {
		total += val
	}
	return total
}
