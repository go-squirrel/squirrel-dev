package res

import "squirrel-dev/internal/pkg/response"

const (
	// 73000-73019: 基础操作
	ErrAppStoreNotFound        = 73001
	ErrDuplicateAppStore       = 73002
	ErrInvalidComposeContent   = 73003
	ErrUnsupportedAppType      = 73004
	ErrInvalidAppStoreConfig  = 73005
	ErrAppStoreUpdateFailed   = 73006
	ErrAppStoreDeleteFailed   = 73007
)

func RegisterCode() {
	// 73000-73019: 基础操作
	response.Register(ErrAppStoreNotFound, "application store not found")
	response.Register(ErrDuplicateAppStore, "application store already exists")
	response.Register(ErrInvalidComposeContent, "invalid compose content")
	response.Register(ErrUnsupportedAppType, "unsupported application type")
	response.Register(ErrInvalidAppStoreConfig, "invalid app store configuration")
	response.Register(ErrAppStoreUpdateFailed, "app store update failed")
	response.Register(ErrAppStoreDeleteFailed, "app store delete failed")
}
