package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrDockerNotInstalled = 10001
	ErrComposeNotFound    = 10002
	ErrComposeStart       = 10003
	ErrComposeCreate      = 10004
	ErrComposeStop        = 10005
)

func RegisterCode() {
	response.Register(ErrDockerNotInstalled, "docker not installed")
	response.Register(ErrComposeNotFound, "docker-compose command not found")
	response.Register(ErrComposeStart, "docker-compose start failed")
	response.Register(ErrComposeCreate, "docker-compose file creation failed")
	response.Register(ErrComposeStop, "docker-compose stop failed")
}
