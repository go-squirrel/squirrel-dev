package deployment

import (
	"fmt"
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// ComposeInfo 存储解析出的 compose 文件信息
type ComposeInfo struct {
	ContainerNames map[string]struct{}
	Ports          map[string]struct{}
	Volumes        map[string]struct{}
	Networks       map[string]struct{}
}

// ServiceConfig compose 文件中的 service 配置
type ServiceConfig struct {
	ContainerName string   `yaml:"container_name"`
	Ports         []string `yaml:"ports"`
	Volumes       []string `yaml:"volumes"`
}

// ComposeConfig docker-compose 文件结构
type ComposeConfig struct {
	Services map[string]ServiceConfig `yaml:"services"`
	Volumes  map[string]any           `yaml:"volumes"`
	Networks map[string]any           `yaml:"networks"`
}

func checkComposeContent(reqContent, modelContent string) (resCode int) {
	if reqContent == "" {
		return res.ErrInvalidDeploymentConfig
	}

	// 如果 modelContent 为空，表示没有已部署的应用，无需检查冲突
	if modelContent == "" {
		return 0
	}

	// 解析要部署的应用的 compose 内容
	reqInfo, err := parseComposeContent(reqContent)
	if err != nil {
		return res.ErrInvalidDeploymentConfig
	}

	// 解析已部署应用的 compose 内容
	modelInfo, err := parseComposeContent(modelContent)
	if err != nil {
		return res.ErrInvalidDeploymentConfig
	}

	// 检查容器名称冲突
	for containerName := range reqInfo.ContainerNames {
		if _, exists := modelInfo.ContainerNames[containerName]; exists {
			return res.ErrComposeContainerNameConflict
		}
	}

	// 检查端口冲突
	for port := range reqInfo.Ports {
		if _, exists := modelInfo.Ports[port]; exists {
			return res.ErrComposePortConflict
		}
	}

	// 检查卷名冲突
	for volume := range reqInfo.Volumes {
		if _, exists := modelInfo.Volumes[volume]; exists {
			return res.ErrComposeVolumeConflict
		}
	}

	// 检查网络名冲突
	for network := range reqInfo.Networks {
		if _, exists := modelInfo.Networks[network]; exists {
			return res.ErrComposeNetworkConflict
		}
	}

	return 0
}

// parseComposeContent 解析 compose 文件内容，提取关键信息
func parseComposeContent(content string) (*ComposeInfo, error) {
	var config ComposeConfig
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse compose content: %w", err)
	}

	info := &ComposeInfo{
		ContainerNames: make(map[string]struct{}),
		Ports:          make(map[string]struct{}),
		Volumes:        make(map[string]struct{}),
		Networks:       make(map[string]struct{}),
	}

	// 提取容器名称和端口
	for _, service := range config.Services {
		// 提取容器名称
		if service.ContainerName != "" {
			info.ContainerNames[service.ContainerName] = struct{}{}
		}

		// 提取端口映射
		for _, portMapping := range service.Ports {
			// 端口映射格式可能是 "host_port:container_port" 或 "host_port:container_port/protocol"
			// 我们只需要主机端口部分
			parts := strings.Split(portMapping, ":")
			if len(parts) > 0 {
				hostPort := parts[0]
				// 提取端口号（去除可能的协议部分）
				if idx := strings.Index(hostPort, "/"); idx != -1 {
					hostPort = hostPort[:idx]
				}
				if hostPort != "" {
					info.Ports[hostPort] = struct{}{}
				}
			}
		}

		// 提取卷名（外部命名卷）
		for _, vol := range service.Volumes {
			// 卷格式可能是 "volume_name:/container/path" 或 "/host/path:/container/path"
			// 只有以命名卷开头的才需要检查冲突
			parts := strings.Split(vol, ":")
			if len(parts) > 0 {
				volumeName := parts[0]
				// 如果是命名卷（不是绝对路径），则添加到卷列表
				if volumeName != "" && !strings.HasPrefix(volumeName, "/") && !strings.HasPrefix(volumeName, "./") && !strings.HasPrefix(volumeName, "~") {
					// 移除可能的读写标志
					if idx := strings.Index(volumeName, ":"); idx != -1 {
						volumeName = volumeName[:idx]
					}
					if volumeName != "" {
						info.Volumes[volumeName] = struct{}{}
					}
				}
			}
		}
	}

	// 提取卷定义
	for volumeName := range config.Volumes {
		info.Volumes[volumeName] = struct{}{}
	}

	// 提取网络定义
	for networkName := range config.Networks {
		info.Networks[networkName] = struct{}{}
	}

	return info, nil
}

// returnDeploymentErrCode 根据错误类型返回精确的部署错误码
func returnDeploymentErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrDeploymentNotFound
	case gorm.ErrDuplicatedKey:
		return res.ErrAlreadyDeployed
	}
	return res.ErrCreateDeploymentRecordFailed
}
