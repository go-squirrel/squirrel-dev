package cache

import "errors"

// 缓存相关错误
var (
	// ErrKeyNotFound key 不存在
	ErrKeyNotFound = errors.New("cache: key not found")
	// ErrCacheNotConnected 缓存未连接
	ErrCacheNotConnected = errors.New("cache: not connected")
	// ErrUnsupportedCacheType 不支持的缓存类型
	ErrUnsupportedCacheType = errors.New("cache: unsupported cache type")
	// ErrInvalidConnectionString 无效的连接字符串
	ErrInvalidConnectionString = errors.New("cache: invalid connection string")
	// ErrNilValue 不能缓存 nil 值
	ErrNilValue = errors.New("cache: cannot cache nil value")
)
