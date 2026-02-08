package res

import "squirrel-dev/internal/pkg/response"

const (
	// 65000-65019: 基础操作
	ErrConfigNotFound          = 65001
	ErrConfigKeyAlreadyExists  = 65002
	ErrInvalidConfigKey        = 65003
	ErrInvalidConfigValue      = 65004
	ErrConfigUpdateFailed      = 65005
	ErrConfigDeleteFailed      = 65006
)

func RegisterCode() {
	// 65000-65019: 基础操作
	response.Register(ErrConfigNotFound, "config not found")
	response.Register(ErrConfigKeyAlreadyExists, "config key already exists")
	response.Register(ErrInvalidConfigKey, "invalid config key")
	response.Register(ErrInvalidConfigValue, "invalid config value")
	response.Register(ErrConfigUpdateFailed, "config update failed")
	response.Register(ErrConfigDeleteFailed, "config delete failed")
}
