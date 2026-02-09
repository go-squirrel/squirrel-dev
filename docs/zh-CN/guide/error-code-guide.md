# Squirrel 错误码规范文档

## 概述

本文档定义了 Squirrel 系统中所有 API 错误码的规范和使用指南。统一的错误码设计有助于：
- 快速定位问题来源
- 区分不同模块的错误
- 便于调试和日志分析
- 提供友好的错误信息

## 错误码设计原则

### 命名规范
- **统一 5 位数**：所有错误码均为 5 位数，便于记忆和识别
- **模块化分段**：每个模块占用独立的 1000 个错误码空间
- **功能细分**：模块内按功能每 20 个错误码为一组
- **预留扩展**：每个模块预留充足空间，满足未来扩展需求

### 错误码结构
```
[模块标识][功能分组][具体错误]
  60      01        01
  ↓        ↓        ↓
  模块    功能组   具体错误
```

## 基础错误码

### 0xxxx: 通用状态
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 0 | `CodeSuccess` | 操作成功 | 所有成功的 API 响应 |

### 41xxx: 通用错误
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 41001 | `ErrCodeParameter` | 参数错误 | 请求参数验证失败 |
| 41002 | `ErrUserOrPassword` | 用户名或密码错误 | 认证失败 |

### 50xxx: 数据库错误
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 50000 | `ErrSQL` | SQL 错误 | 数据库操作异常 |
| 50001 | `ErrSQLNotFound` | 数据未找到 | 数据库查询无结果 |
| 50002 | `ErrSQLNotUnique` | 数据不唯一 | 唯一约束冲突 |
| 50003 | `ErrDuplicatedKey` | 键值重复 | 主键或唯一键冲突 |

**定义位置**: `internal/pkg/response/common.go`

---

## 模块错误码

### 60xxx: 服务器管理 (Server)

**定义位置**: `internal/squ-apiserver/handler/server/res/response_code.go`

#### 60000-60019: 基础操作（增删改查）
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 60001 | `ErrServerNotFound` | 服务器未找到 | 根据查询条件未找到服务器 |
| 60002 | `ErrServerAliasExists` | 服务器别名已存在 | 添加/更新时别名冲突 |
| 60003 | `ErrServerUUIDNotFound` | 通过 UUID 未找到服务器 | Agent 注册时 UUID 不存在 |
| 60004 | `ErrServerAlreadyExists` | 服务器已存在 | 重复添加同一服务器 |
| 60005 | `ErrServerUpdateFailed` | 更新服务器失败 | 更新操作异常 |
| 60006 | `ErrServerDeleteFailed` | 删除服务器失败 | 删除操作异常 |

#### 60020-60039: 验证相关
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 60021 | `ErrInvalidParameter` | 参数验证失败 | 请求参数不符合要求 |
| 60022 | `ErrInvalidAuthType` | 无效的认证类型 | 认证类型不是 password/key/password_key |
| 60023 | `ErrInvalidSSHConfig` | 无效的 SSH 配置 | SSH 配置格式或内容错误 |
| 60024 | `ErrSSHTestFailed` | SSH 连接测试失败 | 无法建立 SSH 连接 |

#### 60040-60059: Agent 通信相关
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 60041 | `ErrConnectFailed` | 连接失败 | 网络连接异常 |
| 60042 | `ErrAgentOffline` | Agent 离线 | Agent 服务不可用 |
| 60043 | `ErrAgentRequestFailed` | Agent 请求失败 | 向 Agent 发送请求失败 |

---

### 70xxx: 应用商店 (App Store)

**定义位置**: `internal/squ-apiserver/handler/app_store/res/response_code.go`

| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 70001 | `ErrAppStoreNotFound` | 应用商店未找到 | 查询的应用商店不存在 |
| 70002 | `ErrDuplicateAppStore` | 应用商店已存在 | 添加重复的应用商店 |
| 70003 | `ErrInvalidComposeContent` | 无效的 Compose 内容 | Docker Compose 格式错误 |
| 70004 | `ErrUnsupportedAppType` | 不支持的应用类型 | 应用类型不在支持列表 |

---

### 71xxx: 应用 (Application)

**定义位置**: `internal/squ-apiserver/handler/application/res/response_code.go`

| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 71001 | `ErrApplicationNotFound` | 应用未找到 | 查询的应用不存在 |
| 71002 | `ErrDuplicateApplication` | 应用已存在 | 添加重复的应用 |
| 71003 | `ErrInvalidAppConfig` | 无效的应用配置 | 应用配置格式或内容错误 |

---

### 72xxx: 部署 (Deployment)

**定义位置**: `internal/squ-apiserver/handler/deployment/res/response_code.go`

#### 72000-72019: 基础操作
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 72001 | `ErrDeploymentNotFound` | 部署记录未找到 | 查询的部署不存在 |
| 72002 | `ErrAlreadyDeployed` | 应用已部署到此服务器 | 重复部署 |
| 72003 | `ErrApplicationNotDeployed` | 应用未部署到此服务器 | 操作未部署的应用 |
| 72004 | `ErrDeployIDGenerateFailed` | 生成部署 ID 失败 | ID 生成器异常 |
| 72005 | `ErrCreateDeploymentRecordFailed` | 创建部署记录失败 | 数据库写入失败 |

#### 72020-72039: Agent 相关
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 72021 | `ErrAgentRequestFailed` | Agent 请求失败 | 发送请求到 Agent 失败 |
| 72022 | `ErrAgentResponseParseFailed` | Agent 响应解析失败 | 响应数据格式错误 |
| 72023 | `ErrAgentDeployFailed` | Agent 部署失败 | Agent 端部署操作失败 |
| 72024 | `ErrAgentDeleteFailed` | Agent 删除应用失败 | Agent 端删除操作失败 |
| 72025 | `ErrAgentStopFailed` | Agent 停止应用失败 | Agent 端停止操作失败 |
| 72026 | `ErrAgentStartFailed` | Agent 启动应用失败 | Agent 端启动操作失败 |
| 72027 | `ErrAgentOperationFailed` | Agent 操作失败 | Agent 端其他操作失败 |

---

### 80xxx: 脚本 (Script)

**定义位置**: `internal/squ-apiserver/handler/script/res/response_code.go`

#### 80000-80019: 基础操作
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 80001 | `ErrScriptNotFound` | 脚本未找到 | 查询的脚本不存在 |
| 80002 | `ErrDuplicateScript` | 脚本已存在 | 添加重复的脚本 |
| 80003 | `ErrInvalidScriptContent` | 无效的脚本内容 | 脚本内容格式错误 |
| 80004 | `ErrUnsupportedScriptType` | 不支持的脚本类型 | 脚本类型不在支持列表 |
| 80005 | `ErrScriptNotDeployed` | 脚本未部署 | 操作未部署的脚本 |

#### 80020-80039: 执行相关
| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 80021 | `ErrScriptExecutionFailed` | 脚本执行失败 | 脚本运行异常 |
| 80022 | `ErrScriptTimeout` | 脚本执行超时 | 脚本运行超过时间限制 |
| 80023 | `ErrServerNotFound` | 服务器未找到（脚本相关） | 脚本关联的服务器不存在 |

---

### 81xxx: 监控 (Monitor)

**定义位置**: `internal/squ-apiserver/handler/monitor/res/response_code.go`

| 错误码 | 常量名 | 描述 | 使用场景 |
|--------|--------|------|----------|
| 81001 | `ErrMonitorFailed` | 监控请求失败 | 获取监控数据失败 |
| 81002 | `ErrInvalidMonitorConfig` | 无效的监控配置 | 监控配置格式错误 |
| 81003 | `ErrMonitorDataNotFound` | 监控数据未找到 | 监控数据不存在 |
| 81004 | `ErrServerNotFound` | 服务器未找到（监控相关） | 监控关联的服务器不存在 |

---

## 错误码分配总览

### 分配规则

| 模块 | 错误码范围 | 说明 |
|------|------------|------|
| 通用状态 | 0 | 成功状态 |
| 通用错误 | 41000-41999 | 参数、认证等通用错误 |
| 数据库错误 | 50000-50999 | SQL 操作相关错误 |
| 服务器管理 | 60000-60999 | Server 模块 |
| 应用商店 | 70000-70999 | App Store 模块 |
| 应用 | 71000-71999 | Application 模块 |
| 部署 | 72000-72999 | Deployment 模块 |
| 脚本 | 80000-80999 | Script 模块 |
| 监控 | 81000-81999 | Monitor 模块 |
| 预留扩展 | 82000-99999 | 未来功能模块 |

### 扩展指南

新增错误码时，请遵循以下步骤：

1. **确认模块范围**：检查对应模块的错误码分配表
2. **选择功能分组**：在模块内找到合适的功能分组（每 20 个为一组）
3. **定义常量**：在对应的 `response_code.go` 中定义常量
4. **注册错误码**：在 `RegisterCode()` 函数中注册错误码和描述
5. **更新文档**：同步更新本文档

**示例**：
```go
// 定义常量
const (
    ErrNewFeatureFailed = 60045 // 新功能失败
)

// 注册错误码
func RegisterCode() {
    response.Register(ErrNewFeatureFailed, "new feature failed")
}
```

---

## 错误码使用示例

### 基础错误返回

```go
func (s *Server) Get(id uint) response.Response {
    daoS, err := s.Repository.Get(id)
    if err != nil {
        return response.Error(res.ErrServerNotFound)
    }
    return response.Success(daoS)
}
```

### 数据库错误处理

```go
func (s *Server) Add(request req.Server) response.Response {
    modelReq := s.requestToModel(request)
    
    err := s.Repository.Add(&modelReq)
    if err != nil {
        // 检查是否是重复键错误
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return response.Error(res.ErrServerAlreadyExists)
        }
        return response.Error(model.ReturnErrCode(err))
    }
    
    return response.Success("success")
}
```

### 参数验证错误

```go
func (s *Server) Update(request req.Server) response.Response {
    if request.Hostname == "" {
        return response.Error(res.ErrInvalidParameter)
    }
    
    modelReq := s.requestToModel(request)
    modelReq.ID = request.ID
    
    err := s.Repository.Update(&modelReq)
    if err != nil {
        return response.Error(res.ErrServerUpdateFailed)
    }
    
    return response.Success("success")
}
```

---

## 注意事项

### 1. 错误码唯一性
- 确保错误码在整个系统中唯一
- 不同模块不能共用同一错误码
- 添加新错误码前先检查是否已存在

### 2. 错误信息友好性
- 错误信息应清晰、准确、易懂
- 使用英文描述，便于国际化
- 避免暴露敏感信息

### 3. 日志记录
- 返回错误码的同时，应记录详细日志
- 日志中应包含错误上下文信息
- 便于问题追踪和调试

### 4. 向后兼容
- 已发布的错误码不得随意更改
- 新增错误码应在文档中标注版本信息
- 废弃错误码应保留至少一个大版本

---

## 附录：相关文件

### 核心文件
- `internal/pkg/response/common.go` - 基础错误码定义
- `internal/pkg/response/response.go` - Response 结构和错误处理函数

### 模块文件
- `internal/squ-apiserver/handler/server/res/response_code.go` - Server 模块
- `internal/squ-apiserver/handler/app_store/res/response_code.go` - App Store 模块
- `internal/squ-apiserver/handler/application/res/response_code.go` - Application 模块
- `internal/squ-apiserver/handler/deployment/res/response_code.go` - Deployment 模块
- `internal/squ-apiserver/handler/script/res/response_code.go` - Script 模块
- `internal/squ-apiserver/handler/monitor/res/response_code.go` - Monitor 模块

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0 | 2026-02-08 | 初始版本，统一错误码规范 |

---

**维护者**: Squirrel 开发团队
**最后更新**: 2026-02-08
