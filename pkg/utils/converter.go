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

	// Check if value exceeds uint range (on 32-bit systems, uint max is math.MaxUint32)
	if value > math.MaxUint && ^uint(0) == math.MaxUint32 {
		return 0, fmt.Errorf("Value %d exceeds uint range", value)
	}

	return uint(value), nil
}
