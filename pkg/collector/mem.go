package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type MemoryStats struct {
	Total     int64 `json:"total"`
	Free      int64 `json:"free"`
	Available int64 `json:"available"`
	Buffers   int64 `json:"buffers"`
	Cached    int64 `json:"cached"`
}

// 返回的数字默认是kb
func GetMemLinuxStats() (stats MemoryStats) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		zap.S().Error("open err:", err)
		return stats
	}
	defer file.Close()

	stats = MemoryStats{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}
		// 直接返回就是kb
		switch fields[0] {
		case "MemTotal:":
			stats.Total = value
		case "MemFree:":
			stats.Free = value
		case "MemAvailable:":
			stats.Available = value
		case "Buffers:":
			stats.Buffers = value
		case "Cached:":
			stats.Cached = value
		}
	}

	return stats
}
