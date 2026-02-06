package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
)

// GenerateServerUUID 生成服务器唯一标识
// 优先使用 /etc/machine-id，其次使用 hostname
func GenerateServerUUID(hostname string) string {
	// 尝试读取 /etc/machine-id（Linux 系统唯一标识）
	machineID, err := readMachineID()
	if err == nil && machineID != "" {
		return machineID
	}

	// 如果没有 machine-id，使用 hostname 生成
	return generateUUIDFromHostname(hostname)
}

// readMachineID 读取系统的 machine-id
func readMachineID() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// generateUUIDFromHostname 从 hostname 生成 UUID
func generateUUIDFromHostname(hostname string) string {
	// 使用 MD5 生成固定长度的 UUID
	hash := md5.Sum([]byte(hostname))
	return hex.EncodeToString(hash[:])
}
