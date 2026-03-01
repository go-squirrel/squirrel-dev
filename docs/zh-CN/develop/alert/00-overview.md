# 告警通知系统设计概览

## 1. 背景

监控模块已实现数据采集和历史查询，但缺少告警触发和通知能力。运维人员需要及时获知系统异常，如：
- CPU/内存使用率过高
- 磁盘空间不足
- 服务不可用
- 自定义指标异常

## 2. 设计目标

- **灵活的告警规则**：支持阈值、持续时间、聚合函数
- **多渠道通知**：邮件、钉钉、企业微信、飞书、Webhook
- **告警管理**：静默、抑制、认领、恢复确认
- **告警历史**：完整的告警生命周期追踪

## 3. 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           告警通知系统架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Monitor Data ──► Rule Engine ──► Alert Manager ──► Notification Channel    │
│        │               │                │                    │              │
│        │               │                │                    ├── Email      │
│        ▼               ▼                ▼                    ├── DingTalk   │
│   ┌─────────┐    ┌──────────┐    ┌───────────┐              ├── WeChat     │
│   │ Metrics │    │ 规则评估  │    │ 告警处理  │              ├── Feishu     │
│   │ CPU/Mem │    │ 阈值判断  │    │ 静默/抑制 │              └── Webhook    │
│   │ Disk... │    │ 聚合计算  │    │ 分组/去重 │                             │
│   └─────────┘    └──────────┘    └───────────┘                             │
│                                        │                                    │
│                                        ▼                                    │
│                                  ┌───────────┐                              │
│                                  │告警历史存储│                              │
│                                  └───────────┘                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 4. 核心功能

| 功能模块 | 说明 |
|---------|------|
| **告警规则** | 定义监控指标、阈值条件、持续时间、告警级别 |
| **规则引擎** | 定时评估规则，触发告警 |
| **告警管理** | 告警状态流转（pending → firing → resolved） |
| **静默机制** | 维护窗口期间暂停告警 |
| **通知渠道** | 支持多种通知方式 |
| **告警历史** | 记录所有告警事件，支持查询和统计 |

## 5. 数据模型概要

```go
// 告警规则
type AlertRule struct {
    ID          uint
    Name        string        // 规则名称
    Metric      string        // 监控指标 (cpu_usage, memory_usage, disk_usage)
    Operator    string        // 比较操作 (>, <, >=, <=, ==)
    Threshold   float64       // 阈值
    Duration    time.Duration // 持续时间
    Severity    string        // 严重级别 (critical, warning, info)
    Labels      map[string]string
    Enabled     bool
}

// 告警事件
type AlertEvent struct {
    ID          uint
    RuleID      uint
    Status      string        // pending, firing, resolved
    Severity    string
    Message     string
    StartsAt    time.Time
    EndsAt      *time.Time
    Labels      map[string]string
    Annotations map[string]string
}

// 通知渠道
type NotificationChannel struct {
    ID       uint
    Name     string
    Type     string        // email, dingtalk, wechat, feishu, webhook
    Config   string        // JSON 配置
    Enabled  bool
}
```

## 6. 开发阶段规划

```
Phase 1: 告警规则引擎
├── 告警规则模型和 CRUD
├── 规则评估器实现
└── 告警事件生成

Phase 2: 告警管理
├── 告警状态流转
├── 静默机制
├── 告警分组和去重
└── 告警历史存储

Phase 3: 通知渠道
├── 通知渠道管理
├── 邮件通知
├── 钉钉/企业微信/飞书
└── Webhook 通知

Phase 4: 前端集成
├── 告警规则管理页面
├── 告警列表页面
├── 通知渠道配置页面
└── 告警详情和操作
```

## 7. 与现有模块的关系

```
Monitor (监控数据) ──► Alert Rule ──► Alert Event
                                            │
                                            ▼
RBAC (权限控制) ◄─── Alert Manager ──► Notification Channel
                                            │
                                            ▼
                                      Audit Log (审计)
```

## 8. 技术选型建议

| 组件 | 建议 | 原因 |
|------|------|------|
| 规则评估 | 自研 + Cron | 简单场景足够，避免引入 Prometheus AlertManager 复杂度 |
| 消息队列 | 可选 Redis | 小规模直接同步发送，大规模可引入队列 |
| 模板引擎 | Go text/template | 支持告警消息模板化 |

## 9. 待详细设计的问题

- [ ] 规则评估频率如何配置？
- [ ] 多服务器聚合告警如何实现？
- [ ] 告警恢复通知是否需要确认？
- [ ] 通知失败重试策略？
- [ ] 告警升级机制（持续未处理升级严重级别）？
