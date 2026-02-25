package cache

import (
	"testing"
)

func TestNew_RistrettoCache(t *testing.T) {
	cache, err := New("memory", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if cache == nil {
		t.Fatal("New() returned nil cache")
	}
	defer cache.Close()

	_, ok := cache.(*RistrettoCache)
	if !ok {
		t.Error("New() did not return RistrettoCache for 'memory' type")
	}
}

func TestNew_RistrettoCacheWithAlias(t *testing.T) {
	cache, err := New("ristretto", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if cache == nil {
		t.Fatal("New() returned nil cache")
	}
	defer cache.Close()
}

func TestNew_UnsupportedType(t *testing.T) {
	cache, err := New("unknown", "")
	if err != ErrUnsupportedCacheType {
		t.Errorf("New() error = %v, want %v", err, ErrUnsupportedCacheType)
	}
	if cache != nil {
		t.Error("New() should return nil for unsupported type")
	}
}

func TestWithOptions(t *testing.T) {
	opts := &Options{}

	WithMaxCost(2 << 30).apply(opts)
	if opts.MaxCost != 2<<30 {
		t.Errorf("WithMaxCost() failed, got %d", opts.MaxCost)
	}

	WithBufferItems(128).apply(opts)
	if opts.BufferItems != 128 {
		t.Errorf("WithBufferItems() failed, got %d", opts.BufferItems)
	}

	WithMetrics(true).apply(opts)
	if !opts.Metrics {
		t.Error("WithMetrics() failed")
	}
}

func TestNew_WithCustomOptions(t *testing.T) {
	cache, err := New("memory", "",
		WithMaxCost(100<<20),
		WithBufferItems(32),
		WithMetrics(true),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer cache.Close()
}
