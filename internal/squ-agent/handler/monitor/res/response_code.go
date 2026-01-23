package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrCodeSystem    = 20001
	ErrCodeParameter = 20002
)

func RegisterCode() {
	response.Register(ErrCodeSystem, "system error")
	response.Register(ErrCodeParameter, "parameter error")
}
