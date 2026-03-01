# 审计日志系统设计概览

## 1. 背景

运维平台涉及敏感操作（服务器管理、应用部署、脚本执行等），需要完整的审计能力：
- 操作追溯：谁在什么时间做了什么操作
- 安全合规：满足安全审计要求
- 问题排查：通过操作历史定位问题

## 2. 设计目标

- **全面记录**：记录所有用户操作和系统事件
- **不可篡改**：审计日志只增不改
- **快速查询**：支持按用户、操作类型、时间范围查询
- **可视化**：操作历史可视化展示

## 3. 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           审计日志系统架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                          Audit Logger                               │    │
│   │  ┌──────────────────────────────────────────────────────────────┐  │    │
│   │  │                     Audit Middleware                          │  │    │
│   │  │  Request ──► Record Request ──► Handler ──► Record Response  │  │    │
│   │  └──────────────────────────────────────────────────────────────┘  │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                          Audit Storage                              │    │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │    │
│   │  │   Database   │  │   Indexer    │  │   Cleaner    │              │    │
│   │  │   日志存储   │  │   索引管理   │  │   过期清理   │              │    │
│   │  └──────────────┘  └──────────────┘  └──────────────┘              │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                          Frontend (Vue)                             │    │
│   │  ├── 审计日志列表                                                   │    │
│   │  ├── 审计日志详情                                                   │    │
│   │  └── 操作统计报表                                                   │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 4. 审计事件类型

| 分类 | 事件 | 说明 |
|------|------|------|
| **认证** | login, logout, login_failed | 登录登出 |
| **服务器** | server_create, server_delete, terminal_connect | 服务器管理 |
| **应用** | app_create, app_deploy, app_start, app_stop | 应用管理 |
| **脚本** | script_create, script_execute, script_delete | 脚本管理 |
| **文件** | file_upload, file_download, file_delete | 文件操作 |
| **用户** | user_create, user_update, user_delete | 用户管理 |
| **角色** | role_create, role_update, role_delete, role_assign | 角色管理 |
| **系统** | config_change, backup_create, restore | 系统配置 |

## 5. 数据模型概要

```go
// 审计日志
type AuditLog struct {
    ID          uint
    Timestamp   time.Time
    
    // 操作者信息
    UserID      uint
    Username    string
    Role        string
    IPAddress   string
    UserAgent   string
    
    // 操作信息
    Action      string        // 操作类型 (server_create, app_deploy...)
    Resource    string        // 资源类型 (server, application, script...)
    ResourceID  uint          // 资源ID
    
    // 请求信息
    Method      string        // HTTP 方法
    Path        string        // 请求路径
    Query       string        // 查询参数
    RequestBody string        // 请求体 (敏感字段脱敏)
    
    // 响应信息
    StatusCode  int           // HTTP 状态码
    Success     bool          // 操作是否成功
    ErrorMsg    string        // 错误信息
    
    // 额外信息
    Duration    int           // 响应时间(毫秒)
    Metadata    string        // 额外元数据 (JSON)
}
```

## 6. 开发阶段规划

```
Phase 1: 审计日志核心
├── 审计日志模型
├── 审计中间件
├── 自动记录 HTTP 请求
└── 日志存储

Phase 2: 日志管理
├── 日志查询 API
├── 时间范围过滤
├── 操作类型过滤
├── 用户过滤
└── 日志详情查看

Phase 3: 前端集成
├── 审计日志列表页面
├── 日志详情弹窗
├── 导出功能
└── 操作统计图表

Phase 4: 增强功能
├── 敏感数据脱敏
├── 日志过期清理
├── 操作统计报表
└── 异常操作告警
```

## 7. 敏感数据脱敏

```go
// 敏感字段列表
var sensitiveFields = []string{
    "password",
    "token",
    "secret",
    "key",
    "credential",
}

// 脱敏处理
func maskSensitiveData(body string) string {
    // 将敏感字段值替换为 ****
}
```

## 8. 与现有模块的关系

```
所有 Handler ──► Audit Middleware ──► AuditLog
                                            │
                    ┌───────────────────────┼───────────────────────┐
                    │                       │                       │
                    ▼                       ▼                       ▼
              Auth (认证)             RBAC (权限)            Alert (告警)
              登录登出审计           权限变更审计           异常操作告警
```

## 9. 配置示例

```yaml
# config/apiserver.yaml
audit:
  enabled: true
  # 排除的路径
  exclude_paths:
    - /api/v1/health
    - /api/v1/metrics
    - /api/v1/monitor/stats
  # 敏感字段脱敏
  sensitive_fields:
    - password
    - token
    - secret
  # 日志保留天数
  retention_days: 90
```

## 10. 待详细设计的问题

- [ ] 日志存储方案？（数据库表 vs 独立存储）
- [ ] 大量日志时如何优化查询性能？
- [ ] 日志是否需要压缩存储？
- [ ] 是否需要支持日志导出和归档？
- [ ] 如何处理 WebSocket 连接的审计？
- [ ] 系统内部操作（定时任务）如何审计？
