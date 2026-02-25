# Monitor Stats 接口性能优化设计文档

## 1. 问题背景

### 1.1 接口描述

`/monitor/stats` 接口用于获取系统监控统计数据，包括 CPU、内存、磁盘、进程等信息。

### 1.2 性能问题

当前接口响应时间过长（约 1-3 秒），主要原因：

| 瓶颈位置 | 具体问题 | 影响 |
|---------|---------|------|
| `pkg/collector/cpu.go:35` | `cpu.Percent(time.Second, true)` 强制等待 1 秒 | 每次请求至少 1 秒延迟 |
| `pkg/collector/process.go:38` | `process.Processes()` 遍历所有系统进程 | 进程数量多时极慢 |
| `pkg/collector/process.go:122` | 每个进程调用 `CPUPercent()` | 大量系统调用 |
| `pkg/collector/process.go:127` | 每个进程调用 `MemoryInfo()` | 大量系统调用 |
| `pkg/collector/factory.go:23` | `CollectAll()` 顺序收集，无并发 | 串行执行累积延迟 |

### 1.3 调用链路

```
GET /monitor/stats
    └── monitor.StatsHandler(service)
        └── service.GetStats()
            ├── factory.CollectAll()           // ~1.1s (CPU 阻塞)
            │   ├── cpu.CollectCPU()           // ~1s (阻塞等待)
            │   ├── memory.CollectMemory()     // ~5ms
            │   └── disk.CollectDisk()         // ~50ms
            └── process.CollectTopCPU(5)       // ~500ms-2s (遍历所有进程)
            └── process.CollectTopMemory(5)    // ~500ms-2s (重复遍历)
```

## 2. 优化目标

| 指标 | 当前 | 目标 |
|------|------|------|
| 接口响应时间 | 1-3s | < 50ms |
| 缓存命中率 | 0% | > 95% |
| 系统负载 | 高（每次请求完整收集） | 低（后台定时收集） |

## 3. 解决方案

### 3.1 核心思路

**缓存 + 异步刷新**：用户请求立即返回缓存数据，后台定时任务异步刷新缓存。

```
┌─────────────────────────────────────────────────────────────┐
│                        请求流程                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  HTTP Request ──► Cache.Get() ──► Immediate Response        │
│                         │                                   │
│                         ▼                                   │
│                    (Cache Hit)                              │
│                         │                                   │
│                    Cache Data                               │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                        后台刷新                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Cron Timer ──► Collector.Collect() ──► Cache.Set()         │
│                      │                                      │
│                      ▼                                      │
│               (异步执行，不阻塞请求)                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 缓存策略

#### 分层 TTL 设计

不同类型的数据变化频率不同，采用分层 TTL 策略：

| 数据类型 | TTL | 理由 |
|---------|-----|------|
| CPU 使用率 | 5s | 变化快，需要较高实时性 |
| 内存使用率 | 10s | 变化较快 |
| 磁盘使用率 | 60s | 变化慢 |
| 进程 TopN | 10s | 变化较快，但收集成本高 |
| 系统信息 | 300s | 基本不变（主机名、CPU型号等） |

#### 缓存 Key 设计

```
monitor:stats:cpu        # CPU 信息
monitor:stats:memory     # 内存信息
monitor:stats:disk       # 磁盘信息
monitor:stats:process    # 进程 TopN 信息
monitor:stats:host       # 主机基础信息
monitor:stats:full       # 完整统计数据（聚合）
```

### 3.3 架构设计

```
┌──────────────────────────────────────────────────────────────────┐
│                           Server                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                         Router                              │  │
│  │   GET /monitor/stats ──► StatsHandler                      │  │
│  └────────────────────────────┬───────────────────────────────┘  │
│                               │                                  │
│  ┌────────────────────────────▼───────────────────────────────┐  │
│  │                    Monitor Service                          │  │
│  │  ┌─────────────────┐  ┌─────────────────┐                  │  │
│  │  │  GetStats()     │  │  Cache          │◄─────────────┐   │  │
│  │  │  - Cache.Get()  │  │  - Ristretto    │              │   │  │
│  │  │  - Fast Return  │  │  - Redis (opt)  │              │   │  │
│  │  └─────────────────┘  └─────────────────┘              │   │  │
│  └─────────────────────────────────────────────────────────┼───┘  │
│                                                            │      │
│  ┌─────────────────────────────────────────────────────────▼───┐  │
│  │                      Cron Scheduler                          │  │
│  │  ┌──────────────────────────────────────────────────────┐   │  │
│  │  │  MonitorStatsJob (每 5s 执行)                         │   │  │
│  │  │  - CollectAll()                                       │   │  │
│  │  │  - Cache.Set()                                        │   │  │
│  │  └──────────────────────────────────────────────────────┘   │  │
│  └──────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
```

## 4. 详细设计

### 4.1 缓存层实现

缓存模块已完成：`internal/pkg/cache/`

```go
// 使用示例
cache, _ := cache.New("memory", "", 
    cache.WithMaxCost(1<<30),  // 1GB
)

// 设置缓存
cache.Set(ctx, "monitor:stats:cpu", cpuInfo, 5*time.Second)

// 获取缓存
val, err := cache.Get(ctx, "monitor:stats:cpu")
```

### 4.2 Monitor Service 改造

**改造前**：每次请求实时收集

```go
func (m *Monitor) GetStats() response.Response {
    hostInfo, err := m.Factory.CollectAll()  // 阻塞 ~1s
    topCPU, _ := procCollector.CollectTopCPU(5)  // 阻塞 ~500ms
    topMemory, _ := procCollector.CollectTopMemory(5)  // 阻塞 ~500ms
    // ...
}
```

**改造后**：优先返回缓存

```go
func (m *Monitor) GetStats() response.Response {
    // 1. 尝试从缓存获取
    if cached, err := m.Cache.Get(ctx, "monitor:stats:full"); err == nil {
        return response.Success(cached)
    }
    
    // 2. 缓存未命中，触发同步收集（降级处理）
    stats := m.collectStats()
    
    // 3. 存入缓存
    m.Cache.Set(ctx, "monitor:stats:full", stats, 10*time.Second)
    
    return response.Success(stats)
}
```

### 4.3 Cron 定时任务

新增定时任务：`internal/squ-agent/cron/monitor_stats.go`

```go
type MonitorStatsJob struct {
    factory *collector.CollectorFactory
    cache   cache.Cache
}

func (j *MonitorStatsJob) Run() {
    ctx := context.Background()
    
    // 收集数据
    hostInfo, _ := j.factory.CollectAll()
    
    // 分层缓存
    j.cache.Set(ctx, "monitor:stats:cpu", hostInfo.CPU, 5*time.Second)
    j.cache.Set(ctx, "monitor:stats:memory", hostInfo.Memory, 10*time.Second)
    j.cache.Set(ctx, "monitor:stats:disk", hostInfo.Disk, 60*time.Second)
    j.cache.Set(ctx, "monitor:stats:full", hostInfo, 10*time.Second)
    
    // 收集进程 TopN
    procCollector := j.factory.GetProcessCollector()
    if procCollector != nil {
        topCPU, _ := procCollector.CollectTopCPU(5)
        topMemory, _ := procCollector.CollectTopMemory(5)
        j.cache.Set(ctx, "monitor:stats:process", ProcessData{
            TopCPU:    topCPU,
            TopMemory: topMemory,
        }, 10*time.Second)
    }
}
```

### 4.4 配置扩展

在 `config/agent.yaml` 添加监控缓存配置：

```yaml
cache:
  type: memory
  monitor:
    enabled: true
    interval: 5s        # 刷新间隔
    cpuTTL: 5s          # CPU 缓存 TTL
    memoryTTL: 10s      # 内存缓存 TTL
    diskTTL: 60s        # 磁盘缓存 TTL
    processTTL: 10s     # 进程缓存 TTL
```

## 5. 实现步骤

### Phase 1: 缓存层集成 (已完成)

- [x] 实现 `internal/pkg/cache` 缓存模块
- [x] 支持 Ristretto 内存缓存
- [x] 支持 Redis 缓存
- [x] 单元测试覆盖

### Phase 2: 配置注入 (已完成)

- [x] 添加缓存配置结构 `internal/squ-agent/config/cache.go`
- [x] 在 `cmd/squ-agent` 初始化缓存
- [x] 更新配置文件 `config/agent.yaml`

### Phase 3: Monitor Service 改造

- [x] 修改 `internal/squ-agent/handler/monitor/service.go`
  - [x] 添加 Cache 依赖注入
- [x] 修改 `internal/squ-agent/handler/monitor/service_stats.go`
  - [x] `GetStats()` 优先返回缓存
  - [x] 缓存未命中时降级处理

### Phase 4: Cron 定时刷新

- [x] 创建 `internal/squ-agent/cron/monitor_stats.go`
  - [x] 实现 `refreshMonitorStatsCache()`
  - [x] 定时收集并存入缓存
- [x] 注册定时任务到 Cron (`cron.go`)
- [x] 修改路由注入 Cache (`router/monitor.go`)

### Phase 5: 测试与验证

- [x] 编译验证
- [ ] 性能基准测试
- [ ] 压力测试

## 6. 性能预期

| 场景 | 改造前 | 改造后 |
|------|--------|--------|
| 缓存命中 | N/A | < 5ms |
| 缓存未命中 | 1-3s | 1-3s（首次请求） |
| 正常请求（缓存命中） | 1-3s | < 50ms |
| 并发 100 请求 | 100-300s 总耗时 | < 1s 总耗时 |

## 7. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 缓存数据过期 | 短暂数据不准确 | 合理设置 TTL，关键场景可强制刷新 |
| 后台任务失败 | 缓存过期 | 添加重试机制，降级到同步收集 |
| 内存占用增加 | 资源消耗 | 设置合理的 MaxCost，监控内存使用 |
| 首次请求慢 | 用户体验 | 服务启动时预热缓存 |

## 8. 后续优化方向

1. **缓存预热**：服务启动时主动收集一次数据
2. **智能 TTL**：根据系统负载动态调整刷新间隔
3. **WebSocket 推送**：实时推送监控数据，无需轮询
4. **指标聚合**：支持历史数据聚合查询
5. **分布式缓存**：多实例场景使用 Redis 共享缓存
