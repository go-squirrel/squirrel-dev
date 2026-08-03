package application

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type composeInfo struct {
	containers map[string]struct{}
	ports      map[string]struct{}
	volumes    map[string]struct{}
	networks   map[string]struct{}
}

type composeService struct {
	ContainerName string   `yaml:"container_name"`
	Ports         []string `yaml:"ports"`
	Volumes       []string `yaml:"volumes"`
}

type composeConfig struct {
	Services map[string]composeService `yaml:"services"`
	Volumes  map[string]any            `yaml:"volumes"`
	Networks map[string]any            `yaml:"networks"`
}

func checkComposeContent(requestContent, deployedContent string) error {
	if requestContent == "" {
		return ErrInvalidConfig
	}
	if deployedContent == "" {
		return nil
	}
	request, err := parseCompose(requestContent)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	deployed, err := parseCompose(deployedContent)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	if intersects(request.containers, deployed.containers) {
		return ErrContainerConflict
	}
	if intersects(request.ports, deployed.ports) {
		return ErrPortConflict
	}
	if intersects(request.volumes, deployed.volumes) {
		return ErrVolumeConflict
	}
	if intersects(request.networks, deployed.networks) {
		return ErrNetworkConflict
	}
	return nil
}

func parseCompose(content string) (*composeInfo, error) {
	var config composeConfig
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return nil, fmt.Errorf("failed to parse compose content: %w", err)
	}
	info := &composeInfo{
		containers: map[string]struct{}{}, ports: map[string]struct{}{},
		volumes: map[string]struct{}{}, networks: map[string]struct{}{},
	}
	for _, service := range config.Services {
		if service.ContainerName != "" {
			info.containers[service.ContainerName] = struct{}{}
		}
		for _, mapping := range service.Ports {
			parts := strings.Split(mapping, ":")
			if len(parts) > 0 {
				port := parts[0]
				if index := strings.Index(port, "/"); index != -1 {
					port = port[:index]
				}
				if port != "" {
					info.ports[port] = struct{}{}
				}
			}
		}
		for _, volume := range service.Volumes {
			parts := strings.Split(volume, ":")
			if len(parts) > 0 {
				name := parts[0]
				if name != "" && !strings.HasPrefix(name, "/") && !strings.HasPrefix(name, "./") && !strings.HasPrefix(name, "~") {
					if index := strings.Index(name, ":"); index != -1 {
						name = name[:index]
					}
					if name != "" {
						info.volumes[name] = struct{}{}
					}
				}
			}
		}
	}
	for name := range config.Volumes {
		info.volumes[name] = struct{}{}
	}
	for name := range config.Networks {
		info.networks[name] = struct{}{}
	}
	return info, nil
}

func intersects(left, right map[string]struct{}) bool {
	for value := range left {
		if _, ok := right[value]; ok {
			return true
		}
	}
	return false
}

func repositoryError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	case gorm.ErrDuplicatedKey:
		return ErrAlreadyDeployed
	default:
		return ErrCreateRecord
	}
}
