package compose

import (
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

// ValidateContent 验证 compose 内容是否有效
func ValidateContent(content string) error {
	// 去除空白字符
	content = strings.TrimSpace(content)

	// 使用 compose-go 解析 compose 文件
	_, err := loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Content: []byte(content),
			},
		},
	})

	return err
}

// IsValidAppType 验证应用类型是否有效
func IsValidAppType(appType string) bool {
	switch appType {
	case "compose", "k8s_manifest", "helm_chart":
		return true
	default:
		return false
	}
}

// TrimSpaceContent 去除内容的空白字符
func TrimSpaceContent(content string) string {
	return strings.TrimSpace(content)
}
