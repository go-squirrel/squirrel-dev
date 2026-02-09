package res

import (
	"squirrel-dev/internal/pkg/response"
)

const (
	// 71000-71019: 基础操作
	ErrApplicationNotFound       = 71001
	ErrDuplicateApplication      = 71002
	ErrInvalidApplicationName    = 71003
	ErrInvalidApplicationType   = 71004
	ErrInvalidApplicationConfig = 71005
	ErrApplicationUpdateFailed  = 71006
	ErrApplicationDeleteFailed  = 71007
)

func RegisterCode() {
	// 71000-71019: 基础操作
	response.Register(ErrApplicationNotFound, "application not found")
	response.Register(ErrDuplicateApplication, "application already exists")
	response.Register(ErrInvalidApplicationName, "invalid application name")
	response.Register(ErrInvalidApplicationType, "invalid application type")
	response.Register(ErrInvalidApplicationConfig, "invalid application configuration")
	response.Register(ErrApplicationUpdateFailed, "application update failed")
	response.Register(ErrApplicationDeleteFailed, "application delete failed")
}
