package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrInvalidComposeContent = 70001
	ErrUnsupportedAppType    = 70002
	ErrAppStoreNotFound     = 70003
	ErrDuplicateAppStore    = 70004
)

func RegisterCode() {
	response.Register(ErrInvalidComposeContent, "invalid compose content")
	response.Register(ErrUnsupportedAppType, "unsupported application type")
	response.Register(ErrAppStoreNotFound, "application store not found")
	response.Register(ErrDuplicateAppStore, "application store already exists")
}
