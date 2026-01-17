package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrDockerNotInstalled   = 10001
	ErrDockerComposeNotFound = 10002
	ErrDockerComposeStart    = 10003
	ErrDockerComposeCreate   = 10004
)

func RegisterCode() {
	response.Register(ErrDockerNotInstalled, "docker not installed")
	response.Register(ErrDockerComposeNotFound, "docker-compose command not found")
	response.Register(ErrDockerComposeStart, "docker-compose start failed")
	response.Register(ErrDockerComposeCreate, "docker-compose file creation failed")
}
