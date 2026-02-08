# 日志规范

## 1. 概述

本项目使用 [zap](https://github.com/uber-go/zap) 作为日志框架，采用结构化日志记录方式，支持日志轮转和分级存储。

## 2. 日志框架

### 2.1 基本配置

日志初始化位于 `internal/pkg/middleware/log/log.go`，支持以下配置：

- **日志级别**: `debug`, `info`, `warn`, `error`, `fatal`
- **日志文件分离**: 
  - Info 级别及以上：写入 info 日志文件
  - Error 级别：同时写入 error 日志文件
- **日志轮转**: 使用 lumberjack 支持
  - `maxSize`: 单个日志文件最大大小（MB）
  - `maxBackups`: 保留的旧日志文件最大数量
  - `maxAge`: 日志文件最大保留天数

### 2.2 日志格式

采用 JSON 格式输出，包含以下标准字段：

```json
{
  "level": "INFO",
  "time": "2006-01-02T15:04:05.000Z07:00",
  "caller": "handler/service.go:123",
  "msg": "log message",
  "field1": "value1",
  "field2": "value2"
}
```

## 3. 日志级别使用规范

### 3.1 Error

**使用场景**: 
- 系统错误、异常情况
- 可能导致服务功能异常的错误
- 数据库操作失败
- 外部服务调用失败

**示例**:

```go
// Failed to send deployment request
zap.L().Error("failed to send deployment request",
    zap.String("url", agentURL),
    zap.Error(err),
)

// Failed to parse agent response
zap.L().Error("failed to parse agent response",
    zap.String("url", agentURL),
    zap.Error(err),
)

// Agent deployment failed
zap.L().Error("agent deployment failed",
    zap.String("url", agentURL),
    zap.Int("code", agentResp.Code),
    zap.String("message", agentResp.Message),
)
```

### 3.2 Warn

**使用场景**:
- 预期外的但可恢复的情况
- 潜在的问题
- 不影响主要功能的失败

**示例**:

```go
// Failed to get server information, but can continue processing other records
zap.L().Warn("failed to get server information",
    zap.Uint("server_id", appServer.ServerID),
    zap.Error(err),
)
```

### 3.3 Info

**使用场景**:
- 重要的业务操作完成
- 状态变更
- 关键流程节点

**示例**:

```go
// Rollback operation
zap.L().Info("rollback: attempting to delete deployed application on agent",
    zap.String("url", agentDeleteURL),
)

// Status update
zap.L().Info("application status updated",
    zap.Uint64("deploy_id", request.DeployID),
    zap.String("status", request.Status),
)
```

## 4. 日志字段规范

### 4.1 必需字段

| 字段名 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| msg | string | 日志消息 | "failed to send deployment request" |

### 4.2 上下文字段

根据业务场景添加相关字段：

#### 4.2.1 标识符字段

```go
zap.Uint("server_id", serverID)
zap.Uint("application_id", applicationID)
zap.Uint64("deploy_id", deployID)
```

#### 4.2.2 网络请求字段

```go
zap.String("url", requestURL)
zap.Int("status", httpCode)
zap.Duration("cost", time.Since(start))
```

#### 4.2.3 错误信息字段

```go
zap.Error(err)  // 自动记录错误信息和堆栈
```

#### 4.2.4 业务状态字段

```go
zap.String("status", "running")
zap.Int("code", agentResp.Code)
zap.String("message", agentResp.Message)
```

#### 4.2.5 HTTP 请求字段（中间件自动记录）

```go
zap.String("method", c.Request.Method)
zap.String("path", path)
zap.String("query", query)
zap.String("ip", c.ClientIP())
zap.String("user-agent", c.Request.UserAgent())
```

### 4.3 字段命名规范

- 使用小写字母和下划线命名：`server_id`, `user_agent`
- 使用有意义的字段名：`application_id` 而不是 `app_id`
- 保持字段名一致性：同一概念在不同场景使用相同的字段名

## 5. 日志记录最佳实践

### 5.1 使用结构化日志

✅ 推荐：
```go
zap.L().Error("failed to send deployment request",
    zap.String("url", agentURL),
    zap.Error(err),
)
```

❌ 不推荐：
```go
zap.L().Error(fmt.Sprintf("failed to send deployment request: %s, url: %s", err.Error(), agentURL))
```

### 5.2 日志消息要清晰

✅ 推荐：
```go
zap.L().Error("failed to call agent to stop application",
    zap.String("url", agentURL),
    zap.Error(err),
)
```

❌ 不推荐：
```go
zap.L().Error("failed",
    zap.String("url", agentURL),
    zap.Error(err),
)
```

### 5.3 记录关键业务信息

在日志中包含足够的上下文信息以便问题追踪：

```go
zap.L().Error("failed to create application server association record",
    zap.Uint("server_id", request.ServerID),
    zap.Uint("application_id", request.ApplicationID),
    zap.Uint64("deploy_id", deployID),
    zap.Error(err),
)
```

### 5.4 使用适当的日志级别

根据错误的严重程度选择正确的日志级别：
- **Error**: 需要立即关注的错误
- **Warn**: 潜在问题，但系统可以继续运行
- **Info**: 关键业务操作和状态变更

### 5.5 避免循环日志

不要在循环中频繁记录相同级别的日志：

```go
// ❌ Not recommended
for _, item := range items {
    zap.L().Info("processing item", zap.Int("id", item.ID))
}

// ✅ Recommended
zap.L().Info("start processing items", zap.Int("count", len(items)))
for _, item := range items {
    // processing logic
}
zap.L().Info("items processing completed", zap.Int("count", len(items)))
```

## 6. HTTP 请求日志

### 6.1 Gin 中间件自动记录

使用 `GinLogger` 中间件自动记录所有 HTTP 请求：

```go
logger.Info(path,
    zap.Int("status", c.Writer.Status()),
    zap.String("method", c.Request.Method),
    zap.String("path", path),
    zap.String("query", query),
    zap.String("ip", c.ClientIP()),
    zap.String("user-agent", c.Request.UserAgent()),
    zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
    zap.Duration("cost", cost),
)
```

### 6.2 Panic 恢复日志

使用 `GinRecovery` 中间件捕获 panic 并记录堆栈信息：

```go
logger.Error("[Recovery from panic]",
    zap.Any("error", err),
    zap.String("request", string(httpRequest)),
    zap.String("stack", string(debug.Stack())),
)
```

## 7. 日志获取方式

在代码中获取 logger 实例：

```go
import "go.uber.org/zap"

// Use global logger (Recommended - Better Performance)
zap.L().Info("log message",
    zap.String("key", "value"),
    zap.Error(err),
)

// Use sugar logger (Concise syntax but lower performance)
zap.S().Warnf("warning: %v", err)
```

### 7.1 L() vs S() 性能对比

| 特性 | `zap.L()` | `zap.S()` |
|------|-----------|-----------|
| 性能 | 高（类型安全） | 低（使用反射） |
| 语法 | 需要类型化字段 | 类似 fmt.Printf |
| 安全性 | 编译时类型检查 | 运行时类型检查 |
| 推荐场景 | 生产环境、高频日志 | 快速原型、非关键路径 |

**性能建议**:
- 优先使用 `zap.L()` 进行结构化日志记录
- 仅在开发、调试或非关键路径使用 `zap.S()`
- 避免在高频调用的热路径中使用 `zap.S()`

## 8. 示例代码

### 8.1 完整的业务方法日志示例

```go
func (a *Deployment) Deploy(request req.DeployApplication) response.Response {
    // 1. 检查应用是否存在
    app, err := a.AppRepo.Get(request.ApplicationID)
    if err != nil {
        return response.Error(res.ErrApplicationNotDeployed)
    }

    // 2. 检查服务器是否存在
    server, err := a.ServerRepo.Get(request.ServerID)
    if err != nil {
        return response.Error(res.ErrApplicationNotDeployed)
    }

    deployID, err := utils.IDGenerate()
    if err != nil {
        zap.L().Error("failed to generate deployment ID",
            zap.Error(err),
        )
        return response.Error(res.ErrDeployIDGenerateFailed)
    }

    // 3. Send deployment request to agent
    agentURL := utils.GenAgentUrl(a.Config.Agent.Http.Scheme,
        server.IpAddress,
        server.AgentPort,
        a.Config.Agent.Http.BaseUrl,
        "application")

    respBody, err := a.HTTPClient.Post(agentURL, agentReq, nil)
    if err != nil {
        zap.L().Error("failed to send deployment request",
            zap.String("url", agentURL),
            zap.Error(err),
        )
        return response.Error(res.ErrAgentRequestFailed)
    }

    // Parse response, check if deployment succeeded
    var agentResp response.Response
    if err := json.Unmarshal(respBody, &agentResp); err != nil {
        zap.L().Error("failed to parse agent response",
            zap.String("url", agentURL),
            zap.Error(err),
        )
        return response.Error(res.ErrAgentResponseParseFailed)
    }

    if agentResp.Code != 0 {
        zap.L().Error("agent deployment failed",
            zap.String("url", agentURL),
            zap.Int("code", agentResp.Code),
            zap.String("message", agentResp.Message),
        )
        return response.Error(res.ErrAgentDeployFailed)
    }

    // 4. Create application server association record
    appServer := model.Deployment{
        ServerID:      request.ServerID,
        ApplicationID: request.ApplicationID,
        DeployID:      deployID,
    }

    err = a.Repository.Add(&appServer)
    if err != nil {
        zap.L().Error("failed to create application server association record",
            zap.Uint("server_id", request.ServerID),
            zap.Uint("application_id", request.ApplicationID),
            zap.Error(err),
        )

        // Rollback
        zap.L().Info("rollback: attempting to delete deployed application on agent",
            zap.String("url", agentDeleteURL),
        )
        _, err = a.HTTPClient.Post(agentDeleteURL, nil, nil)
        if err != nil {
            zap.L().Error("rollback failed: failed to delete application on agent",
                zap.String("url", agentDeleteURL),
                zap.Error(err),
            )
        }

        return response.Error(returnDeploymentErrCode(err))
    }

    return response.Success("deploy success")
}
```

## 9. 注意事项

1. **不要记录敏感信息**: 避免在日志中记录密码、密钥等敏感数据
2. **控制日志量**: 避免在循环中频繁记录日志，使用计数器代替
3. **保持日志格式一致**: 使用相同的字段名和数据类型
4. **及时清理日志**: 配置合理的日志轮转策略，避免磁盘空间耗尽
5. **错误日志要详细**: Error 级别的日志应该包含足够的上下文信息以便定位问题
6. **使用适当的数据类型**: 使用 `zap.Uint()`、`zap.Int()`、`zap.String()` 等类型化字段
7. **优先使用结构化日志**: 尽可能使用 `zap.L()` 进行结构化日志记录，减少使用 `zap.S()`。`zap.L()` 使用类型化字段，性能更好且更安全；`zap.S()` 是 Sugar API，虽然语法简洁但有性能开销，仅在需要快速原型开发或非关键路径时使用

## 10. 相关文件

- 日志实现: `internal/pkg/middleware/log/log.go`
- 业务日志示例: `internal/squ-apiserver/handler/deployment/service.go`
- 日志中间件示例: `internal/squ-apiserver/handler/deployment/handler.go`
