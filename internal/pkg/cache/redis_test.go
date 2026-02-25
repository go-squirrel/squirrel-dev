package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, *RedisCache) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run() error = %v", err)
	}

	cache := NewRedisCache(&Options{})
	err = cache.Connect(mr.Addr())
	if err != nil {
		mr.Close()
		t.Fatalf("Connect() error = %v", err)
	}

	return mr, cache
}

func TestRedisCache_Connect(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run() error = %v", err)
	}
	defer mr.Close()

	cache := NewRedisCache(&Options{})
	err = cache.Connect(mr.Addr())
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer cache.Close()

	if cache.client == nil {
		t.Error("Connect() did not initialize client")
	}
}

func TestRedisCache_ConnectWithSimpleAddress(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run() error = %v", err)
	}
	defer mr.Close()

	cache := NewRedisCache(&Options{})
	// 使用简单地址格式
	err = cache.Connect(mr.Addr())
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer cache.Close()
}

func TestRedisCache_ConnectWithURL(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run() error = %v", err)
	}
	defer mr.Close()

	cache := NewRedisCache(&Options{})
	// 使用 URL 格式
	err = cache.Connect("redis://" + mr.Addr() + "/0")
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer cache.Close()
}

func TestRedisCache_OperationsWithoutConnect(t *testing.T) {
	cache := &RedisCache{}
	ctx := context.Background()

	_, err := cache.Get(ctx, "key")
	if err != ErrCacheNotConnected {
		t.Errorf("Get() error = %v, want %v", err, ErrCacheNotConnected)
	}

	err = cache.Set(ctx, "key", "value", time.Minute)
	if err != ErrCacheNotConnected {
		t.Errorf("Set() error = %v, want %v", err, ErrCacheNotConnected)
	}

	err = cache.Delete(ctx, "key")
	if err != ErrCacheNotConnected {
		t.Errorf("Delete() error = %v, want %v", err, ErrCacheNotConnected)
	}

	_, err = cache.Exists(ctx, "key")
	if err != ErrCacheNotConnected {
		t.Errorf("Exists() error = %v, want %v", err, ErrCacheNotConnected)
	}

	err = cache.Flush(ctx)
	if err != ErrCacheNotConnected {
		t.Errorf("Flush() error = %v, want %v", err, ErrCacheNotConnected)
	}
}

func TestRedisCache_SetAndGet(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()
	defer cache.Close()

	ctx := context.Background()

	err := cache.Set(ctx, "test-key", "test-value", time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	val, err := cache.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if val != "test-value" {
		t.Errorf("Get() = %v, want %v", val, "test-value")
	}
}

func TestRedisCache_GetNotFound(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()
	defer cache.Close()

	ctx := context.Background()

	_, err := cache.Get(ctx, "non-existent-key")
	if err != ErrKeyNotFound {
		t.Errorf("Get() error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRedisCache_Delete(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()
	defer cache.Close()

	ctx := context.Background()

	cache.Set(ctx, "delete-key", "value", time.Minute)

	err := cache.Delete(ctx, "delete-key")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = cache.Get(ctx, "delete-key")
	if err != ErrKeyNotFound {
		t.Errorf("Get() after delete error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRedisCache_Exists(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()
	defer cache.Close()

	ctx := context.Background()

	// Check non-existent
	exists, err := cache.Exists(ctx, "no-key")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() = true for non-existent key")
	}

	// Set value
	cache.Set(ctx, "exists-key", "value", time.Minute)

	// Check exists
	exists, err = cache.Exists(ctx, "exists-key")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() = false for existing key")
	}
}

func TestRedisCache_Flush(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()
	defer cache.Close()

	ctx := context.Background()

	// Set multiple values
	cache.Set(ctx, "key1", "value1", time.Minute)
	cache.Set(ctx, "key2", "value2", time.Minute)

	// Flush
	err := cache.Flush(ctx)
	if err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	// Verify cleared
	_, err = cache.Get(ctx, "key1")
	if err != ErrKeyNotFound {
		t.Errorf("Get() after flush error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRedisCache_GetClient(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()
	defer cache.Close()

	client := cache.GetClient()
	if client == nil {
		t.Error("GetClient() returned nil")
	}

	_, ok := client.(*redis.Client)
	if !ok {
		t.Error("GetClient() did not return *redis.Client")
	}
}

func TestRedisCache_Close(t *testing.T) {
	mr, cache := setupMiniRedis(t)
	defer mr.Close()

	err := cache.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Close should be idempotent
	err = cache.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}

func TestRedisCache_ConnectFailure(t *testing.T) {
	cache := NewRedisCache(&Options{})

	// 尝试连接不存在的地址
	err := cache.Connect("localhost:12345")
	if err == nil {
		cache.Close()
		t.Error("Connect() should fail for non-existent server")
	}
}

func TestNew_RedisCache(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run() error = %v", err)
	}
	defer mr.Close()

	cache, err := New("redis", mr.Addr())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	_, ok := cache.(*RedisCache)
	if !ok {
		t.Error("New() did not return RedisCache for 'redis' type")
	}
}

// 集成测试：需要真实 Redis 服务
// 运行方式: go test -v -run TestRedisCache_Integration -redis-addr=127.0.0.1:6379
func TestRedisCache_Integration(t *testing.T) {
	addr := "127.0.0.1:6379"

	cache, err := New("redis", addr)
	if err != nil {
		t.Skipf("Redis connection failed: %v (skip integration test)", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Test Set/Get
	if err := cache.Set(ctx, "integration-key", "integration-value", time.Minute); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	val, err := cache.Get(ctx, "integration-key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if val != "integration-value" {
		t.Errorf("Get() = %v, want integration-value", val)
	}

	// Test Exists
	exists, err := cache.Exists(ctx, "integration-key")
	if err != nil || !exists {
		t.Errorf("Exists() = %v, %v, want true", exists, err)
	}

	// Test Delete
	if err := cache.Delete(ctx, "integration-key"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	exists, _ = cache.Exists(ctx, "integration-key")
	if exists {
		t.Error("Key should be deleted")
	}

	t.Log("Redis integration test passed")
}
