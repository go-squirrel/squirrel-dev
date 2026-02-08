# Logging Standard

## 1. Overview

This project uses [zap](https://github.com/uber-go/zap) as the logging framework, employing structured logging with support for log rotation and level-based storage.

## 2. Logging Framework

### 2.1 Basic Configuration

Log initialization is located in `internal/pkg/middleware/log/log.go` and supports the following configurations:

- **Log Levels**: `debug`, `info`, `warn`, `error`, `fatal`
- **Log File Separation**:
  - Info level and above: written to info log file
  - Error level: also written to error log file
- **Log Rotation**: supported by lumberjack
  - `maxSize`: maximum size of single log file (MB)
  - `maxBackups`: maximum number of old log files to keep
  - `maxAge`: maximum number of days to retain log files

### 2.2 Log Format

Logs are output in JSON format with the following standard fields:

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

## 3. Log Level Usage Guidelines

### 3.1 Error

**Use Cases**:
- System errors and exceptions
- Errors that may cause service malfunction
- Database operation failures
- External service call failures

**Examples**:

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

**Use Cases**:
- Unexpected but recoverable situations
- Potential issues
- Failures that do not affect main functionality

**Examples**:

```go
// Failed to get server information, but can continue processing other records
zap.L().Warn("failed to get server information",
    zap.Uint("server_id", appServer.ServerID),
    zap.Error(err),
)
```

### 3.3 Info

**Use Cases**:
- Completion of important business operations
- Status changes
- Key process milestones

**Examples**:

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

## 4. Log Field Guidelines

### 4.1 Required Fields

| Field Name | Type | Description | Example |
|------------|------|-------------|---------|
| msg | string | Log message | "failed to send deployment request" |

### 4.2 Context Fields

Add relevant fields based on business scenarios:

#### 4.2.1 Identifier Fields

```go
zap.Uint("server_id", serverID)
zap.Uint("application_id", applicationID)
zap.Uint64("deploy_id", deployID)
```

#### 4.2.2 Network Request Fields

```go
zap.String("url", requestURL)
zap.Int("status", httpCode)
zap.Duration("cost", time.Since(start))
```

#### 4.2.3 Error Information Fields

```go
zap.Error(err)  // Automatically records error message and stack trace
```

#### 4.2.4 Business Status Fields

```go
zap.String("status", "running")
zap.Int("code", agentResp.Code)
zap.String("message", agentResp.Message)
```

#### 4.2.5 HTTP Request Fields (automatically recorded by middleware)

```go
zap.String("method", c.Request.Method)
zap.String("path", path)
zap.String("query", query)
zap.String("ip", c.ClientIP())
zap.String("user-agent", c.Request.UserAgent())
```

### 4.3 Field Naming Guidelines

- Use lowercase letters and underscores: `server_id`, `user_agent`
- Use meaningful field names: `application_id` instead of `app_id`
- Maintain consistency: use the same field names for the same concept across different scenarios

## 5. Logging Best Practices

### 5.1 Use Structured Logging

✅ Recommended:
```go
zap.L().Error("failed to send deployment request",
    zap.String("url", agentURL),
    zap.Error(err),
)
```

❌ Not Recommended:
```go
zap.L().Error(fmt.Sprintf("failed to send deployment request: %s, url: %s", err.Error(), agentURL))
```

### 5.2 Log Messages Should Be Clear

✅ Recommended:
```go
zap.L().Error("failed to call agent to stop application",
    zap.String("url", agentURL),
    zap.Error(err),
)
```

❌ Not Recommended:
```go
zap.L().Error("failed",
    zap.String("url", agentURL),
    zap.Error(err),
)
```

### 5.3 Record Key Business Information

Include sufficient context information in logs for issue tracking:

```go
zap.L().Error("failed to create application server association record",
    zap.Uint("server_id", request.ServerID),
    zap.Uint("application_id", request.ApplicationID),
    zap.Uint64("deploy_id", deployID),
    zap.Error(err),
)
```

### 5.4 Use Appropriate Log Levels

Select the correct log level based on error severity:
- **Error**: Errors requiring immediate attention
- **Warn**: Potential issues, but system can continue running
- **Info**: Key business operations and status changes

### 5.5 Avoid Loop Logging

Do not frequently log the same level in loops:

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

## 6. HTTP Request Logging

### 6.1 Gin Middleware Automatic Logging

Use `GinLogger` middleware to automatically log all HTTP requests:

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

### 6.2 Panic Recovery Logging

Use `GinRecovery` middleware to catch panics and log stack information:

```go
logger.Error("[Recovery from panic]",
    zap.Any("error", err),
    zap.String("request", string(httpRequest)),
    zap.String("stack", string(debug.Stack())),
)
```

## 7. Logger Access

Get logger instances in code:

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

### 7.1 L() vs S() Performance Comparison

| Feature | `zap.L()` | `zap.S()` |
|---------|-----------|-----------|
| Performance | High (type-safe) | Low (uses reflection) |
| Syntax | Requires typed fields | Similar to fmt.Printf |
| Safety | Compile-time type checking | Runtime type checking |
| Recommended Scenarios | Production, high-frequency logging | Quick prototyping, non-critical paths |

**Performance Recommendations**:
- Prioritize `zap.L()` for structured logging
- Use `zap.S()` only for development, debugging, or non-critical paths
- Avoid using `zap.S()` in hot paths with high-frequency calls

## 8. Example Code

### 8.1 Complete Business Method Logging Example

```go
func (a *Deployment) Deploy(request req.DeployApplication) response.Response {
    // 1. Check if application exists
    app, err := a.AppRepo.Get(request.ApplicationID)
    if err != nil {
        return response.Error(res.ErrApplicationNotDeployed)
    }

    // 2. Check if server exists
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

## 9. Important Notes

1. **Do Not Log Sensitive Information**: Avoid logging passwords, keys, or other sensitive data
2. **Control Log Volume**: Avoid frequently logging in loops; use counters instead
3. **Maintain Consistent Log Format**: Use the same field names and data types
4. **Clean Up Logs Promptly**: Configure reasonable log rotation strategies to avoid disk space exhaustion
5. **Error Logs Should Be Detailed**: Error level logs should contain sufficient context information for issue identification
6. **Use Appropriate Data Types**: Use typed fields like `zap.Uint()`, `zap.Int()`, `zap.String()`
7. **Prefer Structured Logging**: Use `zap.L()` for structured logging as much as possible, minimize usage of `zap.S()`. `zap.L()` uses typed fields with better performance and safety; `zap.S()` is the Sugar API with concise syntax but performance overhead, use it only for quick prototyping or non-critical paths

## 10. Related Files

- Log implementation: `internal/pkg/middleware/log/log.go`
- Business logging examples: `internal/squ-apiserver/handler/deployment/service.go`
- Log middleware examples: `internal/squ-apiserver/handler/deployment/handler.go`
