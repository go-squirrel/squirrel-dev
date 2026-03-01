# Monitor 接口时间范围修复方案

## 问题描述

前端监控页面提供了时间范围选择器（1小时 / 6小时 / 24小时 / 7天），但该功能**前后端完全脱节**，切换时间范围后看到的数据完全相同。

### 根本原因

1. **后端无时间过滤**：所有分页接口只接受 `page` 和 `count`，数据库查询仅 `ORDER BY collect_time DESC LIMIT/OFFSET`，没有 `WHERE collect_time >= ?` 过滤
2. **前端未传递参数**：`loadMonitorData()` 硬编码 `page=1, count=100`，完全忽略 `timeRange` 的值
3. **API Server 纯透传**：API Server 只是把路径原样转发给 Agent，没有传递额外参数的机制

## 当前代码分析

### 数据流

```
前端 (Vue) --> API Server (squ-apiserver) --> Agent (squ-agent) --> SQLite
```

### 前端现状

**`front/src/types/monitor.ts:74`**
```typescript
export type TimeRange = '1h' | '6h' | '24h' | '7d'
```

**`front/src/views/Monitor/components/TimeRangeSelector.vue`** — 渲染 4 个按钮，emit `update:modelValue`

**`front/src/views/Monitor/index.vue:141-143`** — watch 了 `timeRange` 但调用时未传递：
```typescript
watch(timeRange, () => {
  loadMonitorData()  // timeRange 未使用
})
```

**`front/src/views/Monitor/index.vue:93-102`** — 硬编码分页参数：
```typescript
const loadMonitorData = async () => {
  const [base, disk, net] = await Promise.all([
    fetchBaseMonitorHistory(selectedServer.value, 1, 100),  // 硬编码
    fetchDiskIOHistory(selectedServer.value, 1, 100),       // 硬编码
    fetchNetIOHistory(selectedServer.value, 1, 100)         // 硬编码
  ])
}
```

**`front/src/api/monitor.ts`** — API 函数无 `timeRange` 参数：
```typescript
export function fetchBaseMonitorHistory(
  serverId: number, page: number = 1, count: number = 100
): Promise<PageData<BaseMonitorRecord>> {
  return get(`/monitor/base/${serverId}/${page}/${count}`)
}
```

### API Server 现状

**`internal/squ-apiserver/router/monitor.go:26-29`** — 路由只有 path 参数：
```
GET /monitor/base/:serverId/:page/:count
GET /monitor/disk/:serverId/:page/:count
GET /monitor/disk-usage/:serverId/:page/:count
GET /monitor/net/:serverId/:page/:count
```

**`internal/squ-apiserver/handler/monitor/service.go:49-51`** — 转发不含时间参数：
```go
func (m *Monitor) GetBaseMonitorPage(serverID uint, page, count int) response.Response {
    path := fmt.Sprintf("monitor/base/%d/%d", page, count)
    return m.callAgent(serverID, path, "get base monitoring data page")
}
```

**`internal/squ-apiserver/handler/monitor/common.go:12`** — 纯代理调用，无法携带 query 参数：
```go
func (m *Monitor) callAgent(serverID uint, path string, description string) response.Response {
    result := m.AgentClient.Get(context.Background(), server, path, ...)
}
```

### Agent 现状

**`internal/squ-agent/repository/monitor/client.go:58-77`** — 无时间过滤：
```go
func (c *Client) GetBaseMonitorPage(page, pageSize int) ([]model.BaseMonitor, int64, error) {
    offset := (page - 1) * pageSize
    err := c.DB.Model(&model.BaseMonitor{}).Count(&total).Error
    err = c.DB.Order("collect_time DESC").Limit(pageSize).Offset(offset).Find(&monitors).Error
    return monitors, total, nil
}
```

## 修复方案

用 `?range=1h` query 参数替代分页，按时间范围返回数据。

### 方案概览

| 时间范围 | 值 | 含义 |
|---------|-----|------|
| 1小时   | `1h`  | `collect_time >= now - 1h` |
| 6小时   | `6h`  | `collect_time >= now - 6h` |
| 24小时  | `24h` | `collect_time >= now - 24h` |
| 7天     | `7d`  | `collect_time >= now - 7d` |

### 需要修改的文件

#### 1. Agent 数据库查询层

**文件**: `internal/squ-agent/repository/monitor/client.go`

为每个 `Get*MonitorPage` 方法新增 `since time.Time` 参数，将分页改为时间范围查询：

```go
func (c *Client) GetBaseMonitorByTimeRange(since time.Time) ([]model.BaseMonitor, error) {
    var monitors []model.BaseMonitor
    err := c.DB.Where("collect_time >= ?", since).
        Order("collect_time ASC").
        Find(&monitors).Error
    return monitors, err
}

func (c *Client) GetDiskIOMonitorByTimeRange(since time.Time) ([]model.DiskIOMonitor, error) {
    var monitors []model.DiskIOMonitor
    err := c.DB.Where("collect_time >= ?", since).
        Order("collect_time ASC").
        Find(&monitors).Error
    return monitors, err
}

func (c *Client) GetNetworkMonitorByTimeRange(since time.Time) ([]model.NetworkMonitor, error) {
    var monitors []model.NetworkMonitor
    err := c.DB.Where("collect_time >= ?", since).
        Order("collect_time ASC").
        Find(&monitors).Error
    return monitors, err
}

func (c *Client) GetDiskUsageMonitorByTimeRange(since time.Time) ([]model.DiskUsageMonitor, error) {
    var monitors []model.DiskUsageMonitor
    err := c.DB.Where("collect_time >= ?", since).
        Order("collect_time ASC").
        Find(&monitors).Error
    return monitors, err
}
```

#### 2. Agent Handler 层

**文件**: `internal/squ-agent/handler/monitor/handler.go`

新增按时间范围查询的 Handler，从 query 参数中读取 `range`：

```go
func BaseMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
    return func(c *gin.Context) {
        timeRange := c.Query("range")
        if timeRange == "" {
            timeRange = "1h"
        }
        resp := service.GetBaseMonitorByRange(timeRange)
        c.JSON(http.StatusOK, resp)
    }
}
```

在 service 层解析 `range` 值为 `time.Duration`：

```go
func parseTimeRange(rangeStr string) (time.Time, error) {
    now := time.Now()
    switch rangeStr {
    case "1h":
        return now.Add(-1 * time.Hour), nil
    case "6h":
        return now.Add(-6 * time.Hour), nil
    case "24h":
        return now.Add(-24 * time.Hour), nil
    case "7d":
        return now.Add(-7 * 24 * time.Hour), nil
    default:
        return time.Time{}, fmt.Errorf("invalid time range: %s", rangeStr)
    }
}
```

#### 3. Agent 路由

**文件**: `internal/squ-agent/router/monitor.go`

新增路由（或替换原有分页路由）：

```go
group.GET("/monitor/base", monitor.BaseMonitorRangeHandler(service))
group.GET("/monitor/disk", monitor.DiskIOMonitorRangeHandler(service))
group.GET("/monitor/disk-usage", monitor.DiskUsageMonitorRangeHandler(service))
group.GET("/monitor/net", monitor.NetworkMonitorRangeHandler(service))
```

#### 4. API Server Service 层

**文件**: `internal/squ-apiserver/handler/monitor/service.go`

新增带 `timeRange` 参数的方法：

```go
func (m *Monitor) GetBaseMonitorByRange(serverID uint, timeRange string) response.Response {
    path := fmt.Sprintf("monitor/base?range=%s", timeRange)
    return m.callAgent(serverID, path, "get base monitor by time range")
}
```

#### 5. API Server Handler 层

**文件**: `internal/squ-apiserver/handler/monitor/handler.go`

新增 Handler，从 query 参数读取 `range` 并透传：

```go
func BaseMonitorRangeHandler(service *Monitor) func(c *gin.Context) {
    return func(c *gin.Context) {
        serverId := c.Param("serverId")
        serverIdUint, err := utils.StringToUint(serverId)
        if err != nil {
            c.JSON(http.StatusOK, response.Error(res.ErrInvalidMonitorConfig))
            return
        }
        timeRange := c.DefaultQuery("range", "1h")
        resp := service.GetBaseMonitorByRange(serverIdUint, timeRange)
        c.JSON(http.StatusOK, resp)
    }
}
```

#### 6. API Server 路由

**文件**: `internal/squ-apiserver/router/monitor.go`

新增路由（或替换原有分页路由）：

```go
group.GET("/monitor/base/:serverId", monitor.BaseMonitorRangeHandler(service))
group.GET("/monitor/disk/:serverId", monitor.DiskIOMonitorRangeHandler(service))
group.GET("/monitor/disk-usage/:serverId", monitor.DiskUsageMonitorRangeHandler(service))
group.GET("/monitor/net/:serverId", monitor.NetworkMonitorRangeHandler(service))
```

#### 7. 前端 API 层

**文件**: `front/src/api/monitor.ts`

```typescript
export function fetchBaseMonitorHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<BaseMonitorRecord[]> {
  return get(`/monitor/base/${serverId}?range=${range}`)
}

export function fetchDiskIOHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<DiskIORecord[]> {
  return get(`/monitor/disk/${serverId}?range=${range}`)
}

export function fetchNetIOHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<NetworkIORecord[]> {
  return get(`/monitor/net/${serverId}?range=${range}`)
}
```

#### 8. 前端页面

**文件**: `front/src/views/Monitor/index.vue`

将 `timeRange` 传入 API 调用：

```typescript
const loadMonitorData = async () => {
  if (!selectedServer.value) return
  loading.value = true
  try {
    const [base, disk, net] = await Promise.all([
      fetchBaseMonitorHistory(selectedServer.value, timeRange.value),
      fetchDiskIOHistory(selectedServer.value, timeRange.value),
      fetchNetIOHistory(selectedServer.value, timeRange.value)
    ])
    // ...
  }
}
```

## 修改文件清单

| # | 文件 | 改动说明 |
|---|------|---------|
| 1 | `internal/squ-agent/repository/monitor/client.go` | 新增 `Get*MonitorByTimeRange(since time.Time)` 方法 |
| 2 | `internal/squ-agent/handler/monitor/handler.go` | 新增 `*RangeHandler`，读取 `?range=` 参数 |
| 3 | `internal/squ-agent/handler/monitor/service.go` | 新增 `GetBaseMonitorByRange` 等方法，解析 range |
| 4 | `internal/squ-agent/router/monitor.go` | 新增/替换 4 条路由 |
| 5 | `internal/squ-apiserver/handler/monitor/handler.go` | 新增 `*RangeHandler`，透传 `?range=` |
| 6 | `internal/squ-apiserver/handler/monitor/service.go` | 新增带 `timeRange` 参数的 service 方法 |
| 7 | `internal/squ-apiserver/handler/monitor/common.go` | `callAgent` 支持携带 query 参数 |
| 8 | `internal/squ-apiserver/router/monitor.go` | 新增/替换 4 条路由 |
| 9 | `front/src/api/monitor.ts` | API 函数改为接受 `TimeRange` 参数 |
| 10 | `front/src/views/Monitor/index.vue` | `loadMonitorData` 传入 `timeRange.value` |

## 注意事项

- 原有分页接口可保留向后兼容，也可直接替换（仅内部使用无外部依赖）
- `collect_time` 字段已有索引（GORM tag 中 `index`），时间范围查询性能可接受
- 7 天数据量可能较大，可考虑在 Agent 端做数据采样（如每 N 条取 1 条）降低传输量
- `parseTimeRange` 应做白名单校验，仅允许 `1h/6h/24h/7d`，防止非法输入
