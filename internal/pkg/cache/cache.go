package cache

import (
	"context"
	"time"
)

// Cache 定义缓存接口
type Cache interface {
	// Connect 连接缓存服务
	Connect(connectionString string) error
	// Close 关闭连接
	Close() error
	// Get 获取缓存值
	Get(ctx context.Context, key string) (any, error)
	// Set 设置缓存值
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	// Exists 检查 key 是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// Flush 清空所有缓存
	Flush(ctx context.Context) error
	// GetClient 获取底层客户端（可选，用于高级操作）
	GetClient() any
}

// Options 缓存配置选项
type Options struct {
	// MaxCost 最大内存成本（ristretto 使用）
	MaxCost int64
	// BufferItems ristretto buffer 大小
	BufferItems int64
	// Metrics 是否启用指标
	Metrics bool
}

// Option 配置选项接口
type Option interface {
	apply(*Options)
}

type funcOption struct {
	f func(*Options)
}

func (fo *funcOption) apply(o *Options) {
	fo.f(o)
}

func newFuncOption(f func(*Options)) *funcOption {
	return &funcOption{f: f}
}

// WithMaxCost 设置最大内存成本
func WithMaxCost(cost int64) Option {
	return newFuncOption(func(o *Options) {
		o.MaxCost = cost
	})
}

// WithBufferItems 设置 buffer 大小
func WithBufferItems(items int64) Option {
	return newFuncOption(func(o *Options) {
		o.BufferItems = items
	})
}

// WithMetrics 启用指标
func WithMetrics(enabled bool) Option {
	return newFuncOption(func(o *Options) {
		o.Metrics = enabled
	})
}

// New 创建缓存实例
func New(cacheType, connectionString string, opts ...Option) (Cache, error) {
	options := &Options{
		MaxCost:     1 << 30, // 1GB
		BufferItems: 64,
		Metrics:     false,
	}

	for _, opt := range opts {
		opt.apply(options)
	}

	var cache Cache

	switch cacheType {
	case "memory", "ristretto":
		cache = NewRistrettoCache(options)
	case "redis":
		cache = NewRedisCache(options)
	default:
		return nil, ErrUnsupportedCacheType
	}

	if err := cache.Connect(connectionString); err != nil {
		return nil, err
	}

	return cache, nil
}
