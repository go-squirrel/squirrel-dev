package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis 缓存实现
type RedisCache struct {
	client  *redis.Client
	options *Options
}

// NewRedisCache 创建 Redis 缓存实例
func NewRedisCache(opts *Options) *RedisCache {
	return &RedisCache{
		options: opts,
	}
}

// Connect 连接 Redis
func (r *RedisCache) Connect(connectionString string) error {
	// 支持两种格式：
	// 1. 简单地址: "localhost:6379"
	// 2. Redis URL: "redis://user:password@localhost:6379/0"
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		// 如果不是 URL 格式，尝试作为简单地址解析
		opt = &redis.Options{
			Addr: connectionString,
		}
	}

	r.client = redis.NewClient(opt)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Ping(ctx).Err()
}

// Close 关闭 Redis 连接
func (r *RedisCache) Close() error {
	if r.client != nil {
		err := r.client.Close()
		r.client = nil // 确保幂等性
		if err != nil && err.Error() != "redis: client is closed" {
			return err
		}
	}
	return nil
}

// Get 获取缓存值
func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	if r.client == nil {
		return nil, ErrCacheNotConnected
	}

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set 设置缓存值
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.client == nil {
		return ErrCacheNotConnected
	}

	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return ErrCacheNotConnected
	}

	return r.client.Del(ctx, key).Err()
}

// Exists 检查 key 是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	if r.client == nil {
		return false, ErrCacheNotConnected
	}

	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

// Flush 清空所有缓存
func (r *RedisCache) Flush(ctx context.Context) error {
	if r.client == nil {
		return ErrCacheNotConnected
	}

	return r.client.FlushDB(ctx).Err()
}

// GetClient 获取底层 Redis 客户端
func (r *RedisCache) GetClient() interface{} {
	return r.client
}
