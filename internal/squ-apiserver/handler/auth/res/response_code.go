package res

import "squirrel-dev/internal/pkg/response"

const (
	// 66000-66019: 认证相关
	ErrAuthFailed         = 66001
	ErrInvalidCredentials = 66002
	ErrTokenGenerateFailed = 66003
	ErrInvalidToken       = 66004
	ErrTokenExpired       = 66005
)

func RegisterCode() {
	// 66000-66019: 认证相关
	response.Register(ErrAuthFailed, "authentication failed")
	response.Register(ErrInvalidCredentials, "invalid username or password")
	response.Register(ErrTokenGenerateFailed, "failed to generate token")
	response.Register(ErrInvalidToken, "invalid token")
	response.Register(ErrTokenExpired, "token expired")
}
