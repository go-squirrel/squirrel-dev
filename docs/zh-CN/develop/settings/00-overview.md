# 系统配置管理设计概览

## 1. 背景

当前系统配置分散在多个配置文件中，缺少统一的管理界面：
- SMTP 配置、钉钉/企业微信 Webhook 等需要修改配置文件
- 无法在运行时动态调整配置
- 缺少配置变更历史

## 2. 设计目标

- **集中管理**：统一管理系统级配置
- **动态更新**：支持运行时更新配置，无需重启
- **安全存储**：敏感配置加密存储
- **配置审计**：记录配置变更历史

## 3. 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          系统配置管理架构                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                          Frontend (Vue)                             │    │
│   │  ├── 通用设置 (系统名称、Logo、时区)                                 │    │
│   │  ├── 安全设置 (密码策略、会话超时)                                   │    │
│   │  ├── 通知设置 (SMTP、钉钉、企业微信)                                 │    │
│   │  └── 存储设置 (日志保留、监控数据保留)                               │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                        Settings Service                             │    │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │    │
│   │  │ CRUD API     │  │ Validation   │  │ Encryption   │              │    │
│   │  │ 配置管理     │  │ 配置校验     │  │ 敏感加密     │              │    │
│   │  └──────────────┘  └──────────────┘  └──────────────┘              │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                          Settings Storage                           │    │
│   │  ┌──────────────────────────────────────────────────────────────┐  │    │
│   │  │  system_settings 表 (key-value 结构)                         │  │    │
│   │  │  - key: 配置项名称                                           │  │    │
│   │  │  - value: 配置值 (JSON 或字符串)                             │  │    │
│   │  │  - encrypted: 是否加密                                       │  │    │
│   │  └──────────────────────────────────────────────────────────────┘  │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 4. 配置项分类

| 分类 | 配置项 | 说明 |
|------|--------|------|
| **通用设置** | system_name, system_logo, timezone | 系统基础信息 |
| **安全设置** | password_policy, session_timeout, mfa_enabled | 安全策略 |
| **通知设置** | smtp_config, dingtalk_webhook, wechat_webhook | 通知渠道配置 |
| **存储设置** | log_retention_days, monitor_retention_days | 数据保留策略 |
| **登录设置** | login_background, login_notice, allowed_ips | 登录配置 |

## 5. 数据模型概要

```go
// 系统设置
type SystemSetting struct {
    ID          uint
    Key         string    // 配置项名称 (smtp.config, security.password_policy)
    Value       string    // 配置值 (JSON 或字符串)
    Encrypted   bool      // 是否加密存储
    Description string    // 配置说明
    UpdatedAt   time.Time
    UpdatedBy   uint      // 更新人ID
}

// 配置变更历史
type SettingHistory struct {
    ID          uint
    SettingKey  string
    OldValue    string
    NewValue    string
    ChangedBy   uint
    ChangedAt   time.Time
}
```

## 6. 配置项详情

### 6.1 安全设置

```json
{
  "password_policy": {
    "min_length": 8,
    "require_uppercase": true,
    "require_lowercase": true,
    "require_number": true,
    "require_special": false,
    "max_age_days": 90,
    "history_count": 5
  },
  "session_timeout": 30,
  "max_login_attempts": 5,
  "lockout_duration": 15,
  "mfa_enabled": false
}
```

### 6.2 通知设置

```json
{
  "smtp": {
    "host": "smtp.example.com",
    "port": 587,
    "username": "noreply@example.com",
    "password": "encrypted:xxxx",
    "from": "Squirrel <noreply@example.com>",
    "use_tls": true
  },
  "dingtalk": {
    "webhook": "https://oapi.dingtalk.com/robot/send?access_token=xxx",
    "secret": "encrypted:xxxx"
  },
  "wechat": {
    "webhook": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"
  }
}
```

### 6.3 存储设置

```json
{
  "log_retention_days": 30,
  "monitor_retention_days": 90,
  "audit_retention_days": 365,
  "max_storage_gb": 100
}
```

## 7. 开发阶段规划

```
Phase 1: 核心框架
├── 配置模型和存储
├── 配置 CRUD API
├── 配置加密/解密
└── 配置缓存

Phase 2: 安全设置
├── 密码策略实现
├── 会话超时控制
├── 登录失败锁定
└── IP 白名单

Phase 3: 通知设置
├── SMTP 配置管理
├── 钉钉/企业微信配置
├── 配置测试接口
└── 通知模板管理

Phase 4: 前端集成
├── 设置管理页面
├── 配置表单组件
├── 配置测试按钮
└── 变更历史查看
```

## 8. 与现有模块的关系

```
Settings (系统配置)
    │
    ├──► Auth (认证) ── 密码策略、会话超时
    │
    ├──► Alert (告警) ── SMTP、钉钉、企业微信配置
    │
    ├──► Audit (审计) ── 审计日志保留策略
    │
    ├──► Log (日志) ── 日志保留策略
    │
    └──► Monitor (监控) ── 监控数据保留策略
```

## 9. API 设计

```
# 获取配置
GET /api/v1/settings/{category}
# 例: GET /api/v1/settings/security

# 更新配置
PUT /api/v1/settings/{category}
# 例: PUT /api/v1/settings/security
# Body: { "session_timeout": 60 }

# 测试通知渠道
POST /api/v1/settings/notification/test
# Body: { "type": "smtp" | "dingtalk" | "wechat" }

# 获取配置历史
GET /api/v1/settings/{category}/history
```

## 10. 待详细设计的问题

- [ ] 配置项如何分组管理？
- [ ] 敏感配置的加密密钥如何管理？
- [ ] 配置变更后如何通知相关组件？
- [ ] 是否需要配置导入/导出功能？
- [ ] 多实例部署时配置如何同步？
- [ ] 配置项的默认值如何管理？
