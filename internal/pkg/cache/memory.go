package cache

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto"
)

// RistrettoCache ristretto 内存缓存实现
type RistrettoCache struct {
	client  *ristretto.Cache
	options *Options
}

// NewRistrettoCache 创建 Ristretto 缓存实例
func NewRistrettoCache(opts *Options) *RistrettoCache {
	return &RistrettoCache{
		options: opts,
	}
}

// Connect 初始化 ristretto 缓存
func (r *RistrettoCache) Connect(_ string) error {
	config := &ristretto.Config{
		NumCounters: 1e7,     // 10M
		MaxCost:     r.options.MaxCost,
		BufferItems: r.options.BufferItems,
		Metrics:     r.options.Metrics,
	}

	cache, err := ristretto.NewCache(config)
	if err != nil {
		return err
	}

	r.client = cache
	return nil
}

// Close 关闭缓存
func (r *RistrettoCache) Close() error {
	if r.client != nil {
		r.client.Close()
	}
	return nil
}

// Get 获取缓存值
func (r *RistrettoCache) Get(ctx context.Context, key string) (interface{}, error) {
	if r.client == nil {
		return nil, ErrCacheNotConnected
	}

	value, found := r.client.Get(key)
	if !found {
		return nil, ErrKeyNotFound
	}

	return value, nil
}

// Set 设置缓存值
func (r *RistrettoCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.client == nil {
		return ErrCacheNotConnected
	}

	// ristretto 使用 int64 的 cost，这里用 1 作为默认值
	// 实际使用时可以根据 value 大小动态计算
	cost := int64(1)
	if ttl > 0 {
		r.client.SetWithTTL(key, value, cost, ttl)
	} else {
		r.client.Set(key, value, cost)
	}

	return nil
}

// Delete 删除缓存
func (r *RistrettoCache) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return ErrCacheNotConnected
	}

	r.client.Del(key)
	return nil
}

// Exists 检查 key 是否存在
func (r *RistrettoCache) Exists(ctx context.Context, key string) (bool, error) {
	if r.client == nil {
		return false, ErrCacheNotConnected
	}

	_, found := r.client.Get(key)
	return found, nil
}

// Flush 清空所有缓存
func (r *RistrettoCache) Flush(ctx context.Context) error {
	if r.client == nil {
		return ErrCacheNotConnected
	}

	r.client.Clear()
	return nil
}

// GetClient 获取底层 ristretto 客户端
func (r *RistrettoCache) GetClient() interface{} {
	return r.client
}
