# 监控模块设计文档

## 概述

监控系统分为两个部分：
- **Overview 页面** - 实时监控数据展示
- **Monitor 页面** - 历史数据趋势分析

## 架构设计

### 整体架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Frontend (Vue)                            │
│  ├── Overview 页面 - 实时数据                                        │
│  └── Monitor 页面 - 历史数据                                         │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        squ-apiserver (代理层)                        │
│  - 根据 serverId 路由请求到对应的 agent                              │
│  - 统一认证和权限控制                                                │
└─────────────────────────────────────────────────────────────────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              ▼                     ▼                     ▼
┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│   squ-agent (主机1)   │ │   squ-agent (主机2)   │ │   squ-agent (主机N)   │
├──────────────────────┤ ├──────────────────────┤ ├──────────────────────┤
│ • 实时数据采集        │ │ • 实时数据采集        │ │ • 实时数据采集        │
│ • 定时数据存储        │ │ • 定时数据存储        │ │ • 定时数据存储        │
│ • 本地数据库(SQLite)  │ │ • 本地数据库(SQLite)  │ │ • 本地数据库(SQLite)  │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘
```

### 数据模型

#### BaseMonitor - 基础监控数据

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| cpu_usage | float64 | CPU 使用率 (%) |
| cpu_per_core | []float64 | 每核使用率 |
| cpu_model | string | CPU 型号 |
| cpu_mhz | float64 | CPU 频率 |
| load1/load5/load15 | float64 | 系统负载 |
| memory_usage | float64 | 内存使用率 (%) |
| memory_total | uint64 | 内存总量 (bytes) |
| memory_used | uint64 | 已用内存 (bytes) |
| memory_available | uint64 | 可用内存 (bytes) |
| swap_total | uint64 | Swap 总量 |
| swap_used | uint64 | Swap 已用 |
| disk_usage | float64 | 磁盘使用率 (%) |
| disk_total | uint64 | 磁盘总量 |
| disk_used | uint64 | 磁盘已用 |
| collect_time | time | 采集时间 |

#### DiskIOMonitor - 磁盘 IO 监控

> ⚠️ **重要说明**：以下字段存储的是**系统启动以来的累计值**，不是瞬时速率。前端绘制趋势图时需要计算相邻数据点的差值除以时间间隔，得到速率（bytes/s）。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| disk_name | string | 磁盘设备名 (sda, sdb...) |
| read_count | uint64 | 累计读取次数 |
| write_count | uint64 | 累计写入次数 |
| read_bytes | uint64 | 累计读取字节数（系统启动以来） |
| write_bytes | uint64 | 累计写入字节数（系统启动以来） |
| read_time | uint64 | 累计读取时间 (ms) |
| write_time | uint64 | 累计写入时间 (ms) |
| collect_time | time | 采集时间 |

**速率计算公式**：
```
读取速率 (bytes/s) = (当前 read_bytes - 上次 read_bytes) / (当前 collect_time - 上次 collect_time)
写入速率 (bytes/s) = (当前 write_bytes - 上次 write_bytes) / (当前 collect_time - 上次 collect_time)
```

#### NetworkMonitor - 网络监控

> ⚠️ **重要说明**：以下字段存储的是**系统启动以来的累计值**，不是瞬时速率。前端绘制趋势图时需要计算相邻数据点的差值除以时间间隔，得到速率（bytes/s）。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| interface_name | string | 网卡名 (eth0, eth1...) |
| bytes_sent | uint64 | 累计发送字节数（系统启动以来） |
| bytes_recv | uint64 | 累计接收字节数（系统启动以来） |
| packets_sent | uint64 | 累计发送包数 |
| packets_recv | uint64 | 累计接收包数 |
| err_in | uint64 | 累计接收错误数 |
| err_out | uint64 | 累计发送错误数 |
| drop_in | uint64 | 累计接收丢包数 |
| drop_out | uint64 | 累计发送丢包数 |
| collect_time | time | 采集时间 |

**速率计算公式**：
```
上传速率 (bytes/s) = (当前 bytes_sent - 上次 bytes_sent) / (当前 collect_time - 上次 collect_time)
下载速率 (bytes/s) = (当前 bytes_recv - 上次 bytes_recv) / (当前 collect_time - 上次 collect_time)
```

---

## 前端页面设计

### 页面定位

| 页面 | 路由 | 数据类型 | 用途 |
|------|------|---------|------|
| Overview | `/` | 实时数据 | 快速概览当前服务器状态 |
| Monitor | `/monitor` | 历史数据 | 深度分析历史趋势 |

### Monitor 页面布局

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  📊 监控中心                                        [服务器选择 ▼]          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────┐  ┌─────────────────────────┐                  │
│  │      📈 表盘1: CPU       │  │     📈 表盘2: 内存       │                  │
│  │  ─────────────────────  │  │  ─────────────────────  │                  │
│  │  │ ▁▂▃▄▅▆▇█▇▆▅▄▃▂▁   │  │  │ ▁▂▃▄▅▆▇█▇▆▅▄▃▂▁   │                  │
│  │  │   历史趋势图表        │  │  │   历史趋势图表        │                  │
│  │  ─────────────────────  │  │  ─────────────────────  │                  │
│  │  数据: cpu_usage        │  │  数据: memory_usage     │                  │
│  └─────────────────────────┘  └─────────────────────────┘                  │
│                                                                             │
│  ┌─────────────────────────┐  ┌─────────────────────────┐                  │
│  │   💾 表盘3: 磁盘IO       │  │   🌐 表盘4: 网络IO       │                  │
│  │        [sda ▼]          │  │        [eth0 ▼]         │                  │
│  │  ─────────────────────  │  │  ─────────────────────  │                  │
│  │  │ ▁▂▃▄▅▆▇█▇▆▅▄▃▂▁   │  │  │ ▁▂▃▄▅▆▇█▇▆▅▄▃▂▁   │                  │
│  │  │   历史趋势图表        │  │  │   历史趋势图表        │                  │
│  │  ─────────────────────  │  │  ─────────────────────  │                  │
│  │  默认: 所有磁盘总量      │  │  默认: 所有网卡总量      │                  │
│  │  可选: sda, sdb, all    │  │  可选: eth0, eth1, all  │                  │
│  └─────────────────────────┘  └─────────────────────────┘                  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  ⏱️ 时间范围      [1小时] [6小时] [24小时] [7天] [自定义]               ││
│  └─────────────────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
```

### 表盘设计说明

| 表盘 | 数据来源 | 默认展示 | 筛选选项 | 说明 |
|------|---------|---------|---------|------|
| **表盘1** | BaseMonitor.cpu_usage | CPU使用率趋势 | 无 | 单一指标，无需筛选 |
| **表盘2** | BaseMonitor.memory_usage | 内存使用率趋势 | 无 | 单一指标，无需筛选 |
| **表盘3** | DiskIOMonitor | 所有磁盘IO总量 | sda, sdb, ... , all | 支持选择单个磁盘或总量 |
| **表盘4** | NetworkMonitor | 所有网卡流量总量 | eth0, eth1, ... , all | 支持选择单个网卡或总量 |

### 交互设计

#### 1. 服务器选择器
- 位置：页面右上角
- 功能：切换查看不同服务器的监控数据
- 触发：选择后刷新所有表盘数据

#### 2. 时间范围选择器
- 位置：页面底部
- 选项：1小时、6小时、24小时、7天、自定义
- 功能：统一控制所有表盘的时间范围

#### 3. 磁盘/网卡选择器
- 位置：各自表盘右上角
- 功能：筛选特定设备的数据
- 默认值：`all` (显示总量)

---

## 后端 API 设计

### 现有 API (实时数据)

```
# 获取服务器实时监控统计
GET /api/v1/monitor/stats/{serverId}

响应:
{
  "cpu_usage": 45.2,
  "cpu_per_core": [40.1, 50.3, ...],
  "memory_usage": 67.5,
  "disk_usage": 58.0,
  ...
}
```

### 需要扩展的 API (历史数据)

```
# 获取基础监控历史数据
GET /api/v1/monitor/base/history/{serverId}?page=1&count=100

响应:
{
  "list": [
    {
      "id": 1,
      "cpu_usage": 45.2,
      "memory_usage": 67.5,
      "disk_usage": 58.0,
      "collect_time": "2024-01-15T10:00:00Z"
    },
    ...
  ],
  "total": 1000,
  "page": 1,
  "size": 100
}

# 获取磁盘IO历史数据
GET /api/v1/monitor/diskio/history/{serverId}?device=sda&page=1&count=100
# device 参数可选，不传则返回所有磁盘数据（需前端聚合）

响应:
{
  "list": [
    {
      "id": 1,
      "disk_name": "sda",
      "read_bytes": 120000000,    // 累计值，非速率
      "write_bytes": 45000000,    // 累计值，非速率
      "collect_time": "2024-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "disk_name": "sda",
      "read_bytes": 120500000,    // 相比上一条增加了 500KB
      "write_bytes": 45100000,    // 相比上一条增加了 100KB
      "collect_time": "2024-01-15T10:01:00Z"
    },
    ...
  ]
}
# 前端计算速率示例：
# read_speed = (120500000 - 120000000) / 60 = 8333 bytes/s ≈ 8KB/s
# write_speed = (45100000 - 45000000) / 60 = 1666 bytes/s ≈ 1.6KB/s

# 获取网络IO历史数据
GET /api/v1/monitor/netio/history/{serverId}?interface=eth0&page=1&count=100
# interface 参数可选，不传则返回所有网卡数据

响应:
{
  "list": [
    {
      "id": 1,
      "interface_name": "eth0",
      "bytes_sent": 1200000000,   // 累计值，非速率
      "bytes_recv": 3400000000,   // 累计值，非速率
      "collect_time": "2024-01-15T10:00:00Z"
    },
    ...
  ]
}

# 获取设备列表（用于下拉选择）
GET /api/v1/monitor/devices/{serverId}

响应:
{
  "disks": ["sda", "sdb", "sdc"],
  "interfaces": ["eth0", "eth1"]
}
```

---

## 前端实现计划

### 文件结构

```
front/src/
├── views/
│   └── Monitor/
│       ├── index.vue                    # 监控中心主页
│       ├── components/
│       │   ├── CPUMonitorChart.vue      # CPU历史趋势图表
│       │   ├── MemoryMonitorChart.vue   # 内存历史趋势图表
│       │   ├── DiskIOChart.vue          # 磁盘IO图表(含筛选)
│       │   ├── NetIOChart.vue           # 网络IO图表(含筛选)
│       │   └── TimeRangeSelector.vue    # 时间范围选择器
│       └── composables/
│           ├── useMonitorHistory.ts     # 历史数据查询逻辑
│           └── useDeviceList.ts         # 设备列表获取逻辑
├── api/
│   └── monitor.ts                       # 监控相关API (扩展)
└── types/
    └── monitor.ts                       # 监控类型定义 (扩展)
```

### 类型定义

```typescript
// front/src/types/monitor.ts

// 基础监控历史记录
export interface BaseMonitorRecord {
  id: number
  cpu_usage: number
  cpu_per_core: number[]
  cpu_model: string
  cpu_mhz: number
  load1: number
  load5: number
  load15: number
  memory_usage: number
  memory_total: number
  memory_used: number
  memory_available: number
  swap_total: number
  swap_used: number
  disk_usage: number
  disk_total: number
  disk_used: number
  collect_time: string
}

// 磁盘IO历史记录（累计值）
export interface DiskIORecord {
  id: number
  disk_name: string
  read_count: number
  write_count: number
  read_bytes: number      // 累计值，需计算速率
  write_bytes: number     // 累计值，需计算速率
  read_time: number
  write_time: number
  collect_time: string
}

// 磁盘IO速率记录（前端计算后）
export interface DiskIOSpeedRecord {
  collect_time: string
  read_speed: number      // bytes/s
  write_speed: number     // bytes/s
}

// 网络IO历史记录（累计值）
export interface NetworkIORecord {
  id: number
  interface_name: string
  bytes_sent: number      // 累计值，需计算速率
  bytes_recv: number      // 累计值，需计算速率
  packets_sent: number
  packets_recv: number
  err_in: number
  err_out: number
  drop_in: number
  drop_out: number
  collect_time: string
}

// 网络IO速率记录（前端计算后）
export interface NetworkIOSpeedRecord {
  collect_time: string
  upload_speed: number    // bytes/s
  download_speed: number  // bytes/s
}

// 设备列表
export interface DeviceList {
  disks: string[]
  interfaces: string[]
}

// 分页数据
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  size: number
}
```

### API 扩展

```typescript
// front/src/api/monitor.ts

import { get } from '@/utils/request'
import type { PageData, BaseMonitorRecord, DiskIORecord, NetworkIORecord, DeviceList } from '@/types/monitor'

// 获取基础监控历史
export function fetchBaseMonitorHistory(
  serverId: number,
  params: { page: number; count: number }
): Promise<PageData<BaseMonitorRecord>> {
  return get(`/monitor/base/history/${serverId}`, { params })
}

// 获取磁盘IO历史
export function fetchDiskIOHistory(
  serverId: number,
  params: { page: number; count: number; device?: string }
): Promise<PageData<DiskIORecord>> {
  return get(`/monitor/diskio/history/${serverId}`, { params })
}

// 获取网络IO历史
export function fetchNetIOHistory(
  serverId: number,
  params: { page: number; count: number; interface?: string }
): Promise<PageData<NetworkIORecord>> {
  return get(`/monitor/netio/history/${serverId}`, { params })
}

// 获取设备列表
export function fetchDeviceList(serverId: number): Promise<DeviceList> {
  return get(`/monitor/devices/${serverId}`)
}
```

### 速率计算工具函数

```typescript
// front/src/utils/monitor.ts

import type { DiskIORecord, NetworkIORecord, DiskIOSpeedRecord, NetworkIOSpeedRecord } from '@/types/monitor'

/**
 * 将磁盘IO累计值转换为速率
 * @param records 按时间正序排列的历史记录
 * @returns 速率记录数组
 */
export function calculateDiskIOSpeed(records: DiskIORecord[]): DiskIOSpeedRecord[] {
  if (records.length < 2) return []
  
  const result: DiskIOSpeedRecord[] = []
  
  for (let i = 1; i < records.length; i++) {
    const prev = records[i - 1]
    const curr = records[i]
    
    // 计算时间差（秒）
    const timeDiff = (new Date(curr.collect_time).getTime() - new Date(prev.collect_time).getTime()) / 1000
    
    if (timeDiff <= 0) continue
    
    // 计算速率（bytes/s）
    result.push({
      collect_time: curr.collect_time,
      read_speed: Math.max(0, (curr.read_bytes - prev.read_bytes) / timeDiff),
      write_speed: Math.max(0, (curr.write_bytes - prev.write_bytes) / timeDiff),
    })
  }
  
  return result
}

/**
 * 将网络IO累计值转换为速率
 * @param records 按时间正序排列的历史记录
 * @returns 速率记录数组
 */
export function calculateNetworkIOSpeed(records: NetworkIORecord[]): NetworkIOSpeedRecord[] {
  if (records.length < 2) return []
  
  const result: NetworkIOSpeedRecord[] = []
  
  for (let i = 1; i < records.length; i++) {
    const prev = records[i - 1]
    const curr = records[i]
    
    // 计算时间差（秒）
    const timeDiff = (new Date(curr.collect_time).getTime() - new Date(prev.collect_time).getTime()) / 1000
    
    if (timeDiff <= 0) continue
    
    // 计算速率（bytes/s）
    result.push({
      collect_time: curr.collect_time,
      upload_speed: Math.max(0, (curr.bytes_sent - prev.bytes_sent) / timeDiff),
      download_speed: Math.max(0, (curr.bytes_recv - prev.bytes_recv) / timeDiff),
    })
  }
  
  return result
}
```

> 💡 **注意**：上述计算函数要求输入数据按时间**正序**排列（旧→新）。从后端获取的数据通常是倒序（新→旧），需要先反转数组。

---

## 后端实现要点

### 数据采集 (squ-agent)

定时任务已实现，参考 `internal/squ-agent/cron/monitor.go`:
- 默认每 60 秒采集一次
- 数据存储到本地 SQLite 数据库

### 历史数据查询 (squ-agent)

需要在 `internal/squ-agent/handler/monitor/` 扩展:

```go
// service_history.go

// GetBaseMonitorHistory 获取基础监控历史数据
func (s *MonitorService) GetBaseMonitorHistory(page, count int) (*PageData[model.BaseMonitor], error) {
    var records []model.BaseMonitor
    var total int64
    
    db.DB.Model(&model.BaseMonitor{}).Count(&total)
    db.DB.Order("collect_time DESC").
        Offset((page - 1) * count).
        Limit(count).
        Find(&records)
    
    return &PageData[model.BaseMonitor]{
        List:  records,
        Total: total,
        Page:  page,
        Size:  count,
    }, nil
}

// GetDiskIOHistory 获取磁盘IO历史
func (s *MonitorService) GetDiskIOHistory(page, count int, device string) (*PageData[model.DiskIOMonitor], error) {
    query := db.DB.Model(&model.DiskIOMonitor{})
    if device != "" && device != "all" {
        query = query.Where("disk_name = ?", device)
    }
    // ...
}

// GetNetIOHistory 获取网络IO历史
func (s *MonitorService) GetNetIOHistory(page, count int, iface string) (*PageData[model.NetworkMonitor], error) {
    query := db.DB.Model(&model.NetworkMonitor{})
    if iface != "" && iface != "all" {
        query = query.Where("interface_name = ?", iface)
    }
    // ...
}

// GetDeviceList 获取设备列表
func (s *MonitorService) GetDeviceList() (*DeviceList, error) {
    var disks []string
    var interfaces []string
    
    db.DB.Model(&model.DiskIOMonitor{}).
        Distinct("disk_name").
        Pluck("disk_name", &disks)
    
    db.DB.Model(&model.NetworkMonitor{}).
        Distinct("interface_name").
        Pluck("interface_name", &interfaces)
    
    return &DeviceList{
        Disks:      disks,
        Interfaces: interfaces,
    }, nil
}
```

---

## 开发优先级

1. **P0 - 核心功能**
   - [ ] 后端：扩展历史数据查询 API
   - [ ] 前端：创建 Monitor 页面路由
   - [ ] 前端：实现 CPU/内存趋势图表

2. **P1 - 增强功能**
   - [ ] 前端：磁盘IO图表 + 设备筛选
   - [ ] 前端：网络IO图表 + 设备筛选
   - [ ] 前端：时间范围选择器

3. **P2 - 优化体验**
   - [ ] 图表数据聚合（按小时/天）
   - [ ] 图表交互（缩放、tooltip）
   - [ ] 数据缓存优化

---

## 相关文件

- 后端 Agent 处理器: `internal/squ-agent/handler/monitor/`
- 后端 API Server 代理: `internal/squ-apiserver/handler/monitor/`
- 后端数据模型: `internal/squ-agent/model/monitor.go`
- 后端定时任务: `internal/squ-agent/cron/monitor.go`
- 前端 API: `front/src/api/monitor.ts`
- 前端 Overview: `front/src/views/Overview/`
