package compose

import (
	"context"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

// ValidateContent 验证 compose 内容是否有效
func ValidateContent(projectName string, content string) error {
	// 去除空白字符
	content = strings.TrimSpace(content)
	projectName = strings.ToLower(projectName)
	// 使用 compose-go 解析 compose 文件
	_, err := loader.LoadWithContext(
		context.Background(),
		types.ConfigDetails{
			ConfigFiles: []types.ConfigFile{
				{
					Content: []byte(content),
				},
			},
		},
		func(options *loader.Options) {
			options.SetProjectName(projectName, true) // true = 忽略 YAML 中的 name（如果有的话）
		},
	)
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
