package config

// Cache 缓存配置
type Cache struct {
	// Type 缓存类型: memory(默认) 或 redis
	Type string `mapstructure:"type"`
	// Redis Redis 配置
	Redis RedisCacheConfig `mapstructure:"redis"`
	// Memory 内存缓存配置
	Memory MemoryCacheConfig `mapstructure:"memory"`
}

// RedisCacheConfig Redis 缓存配置
type RedisCacheConfig struct {
	// Addr Redis 地址，格式: host:port
	Addr string `mapstructure:"addr"`
	// Password Redis 密码
	Password string `mapstructure:"password"`
	// DB Redis 数据库编号
	DB int `mapstructure:"db"`
	// PoolSize 连接池大小
	PoolSize int `mapstructure:"poolSize"`
}

// MemoryCacheConfig 内存缓存配置 (ristretto)
type MemoryCacheConfig struct {
	// MaxCost 最大内存成本（字节），默认 1GB
	MaxCost int64 `mapstructure:"maxCost"`
	// BufferItems buffer 大小，默认 64
	BufferItems int64 `mapstructure:"bufferItems"`
	// Metrics 是否启用指标
	Metrics bool `mapstructure:"metrics"`
}

// GetConnectionString 获取 Redis 连接字符串
func (c *Cache) GetConnectionString() string {
	if c.Type == "redis" && c.Redis.Addr != "" {
		return c.Redis.Addr
	}
	return ""
}
