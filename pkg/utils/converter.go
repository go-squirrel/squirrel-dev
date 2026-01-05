package utils

import (
	"fmt"
	"math"
	"strconv"
)

func StringToUint(s string) (uint, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	// 检查是否超出 uint 范围（在32位系统上 uint 最大为 math.MaxUint32）
	if value > math.MaxUint && ^uint(0) == math.MaxUint32 {
		return 0, fmt.Errorf("值 %d 超出 uint 范围", value)
	}

	return uint(value), nil
}
