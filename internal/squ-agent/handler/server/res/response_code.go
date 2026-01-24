package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrCodeSystem    = 30001
	ErrCodeParameter = 30002
)

func RegisterCode() {
	response.Register(ErrCodeSystem, "system error")
	response.Register(ErrCodeParameter, "parameter error")
}
