package cache

import (
	"context"
	"testing"
	"time"
)

func TestRistrettoCache_Connect(t *testing.T) {
	cache := NewRistrettoCache(&Options{
		MaxCost:     1 << 20,
		BufferItems: 64,
	})

	err := cache.Connect("")
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer cache.Close()

	if cache.client == nil {
		t.Error("Connect() did not initialize client")
	}
}

func TestRistrettoCache_OperationsWithoutConnect(t *testing.T) {
	cache := &RistrettoCache{}
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

func TestRistrettoCache_SetAndGet(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Set value
	err = cache.Set(ctx, "test-key", "test-value", time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Ristretto 异步写入，需要等待
	time.Sleep(10 * time.Millisecond)

	// Get value
	val, err := cache.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if val != "test-value" {
		t.Errorf("Get() = %v, want %v", val, "test-value")
	}
}

func TestRistrettoCache_GetNotFound(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	_, err = cache.Get(ctx, "non-existent-key")
	if err != ErrKeyNotFound {
		t.Errorf("Get() error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRistrettoCache_Delete(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Set value
	cache.Set(ctx, "delete-key", "value", time.Minute)
	time.Sleep(10 * time.Millisecond)

	// Delete
	err = cache.Delete(ctx, "delete-key")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deleted
	_, err = cache.Get(ctx, "delete-key")
	if err != ErrKeyNotFound {
		t.Errorf("Get() after delete error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRistrettoCache_Exists(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
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
	time.Sleep(10 * time.Millisecond)

	// Check exists
	exists, err = cache.Exists(ctx, "exists-key")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() = false for existing key")
	}
}

func TestRistrettoCache_Flush(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Set multiple values
	cache.Set(ctx, "key1", "value1", time.Minute)
	cache.Set(ctx, "key2", "value2", time.Minute)
	time.Sleep(10 * time.Millisecond)

	// Flush
	err = cache.Flush(ctx)
	if err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	// Verify cleared
	_, err = cache.Get(ctx, "key1")
	if err != ErrKeyNotFound {
		t.Errorf("Get() after flush error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRistrettoCache_TTL(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	ctx := context.Background()

	// Set with short TTL
	cache.Set(ctx, "ttl-key", "value", 50*time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	// Should exist immediately
	exists, _ := cache.Exists(ctx, "ttl-key")
	if !exists {
		t.Error("Key should exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Should be expired
	_, err = cache.Get(ctx, "ttl-key")
	if err != ErrKeyNotFound {
		t.Errorf("Get() after TTL error = %v, want %v", err, ErrKeyNotFound)
	}
}

func TestRistrettoCache_GetClient(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()

	client := cache.GetClient()
	if client == nil {
		t.Error("GetClient() returned nil")
	}
}

func TestRistrettoCache_Close(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = cache.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Close should be idempotent
	err = cache.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}
