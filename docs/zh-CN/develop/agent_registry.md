# Agent 注册机制设计方案

## 背景

当前 `squ-agent` 加入集群的操作存在以下问题：

1. **UUID 生成位置不对**：UUID 在 apiserver 端生成，但 agent 无法知道自己被分配的 UUID
2. **没有"加入集群"的完整流程**：agent 没有主动注册的逻辑
3. **缺少防重复加入机制**：没有检查 agent 是否已经加入其他集群

## 当前代码分析

### APIServer 端

- `Add` 方法：通过前端页面添加服务器，生成 UUID
- `Registry` 方法：已实现但**没有暴露为路由**
- Server model 有 `UUID` 字段作为唯一标识

### Agent 端

- 配置文件 `agent.yaml` 中有 `apiserver.http` 配置指向 apiserver
- `registerAgent()` 方法已定义但**只是打印日志，没有实际实现**
- agent 本地存储了 `server_id` 配置

## 设计方案

### 整体流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           方式一：从 APIServer 添加                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  用户在页面填写服务器信息 ──► APIServer 生成 UUID ──► 保存记录(状态: pending) │
│                                                                             │
│                 用户在目标服务器安装 agent，配置中填入 UUID                   │
│                              ▼                                              │
│                agent 启动时携带 UUID 向 APIServer 注册                       │
│                              ▼                                              │
│           APIServer 验证 UUID ──► 更新服务器状态为 active                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                           方式二：Agent 主动加入                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│       用户在目标服务器安装 agent，配置中填入 APIServer 地址                   │
│                              ▼                                              │
│                   agent 启动时生成自己的 UUID                                │
│                              ▼                                              │
│                  agent 向 APIServer 发起注册请求                             │
│                              ▼                                              │
│     ┌────────────────────────┼────────────────────────┐                     │
│     ▼                        ▼                        ▼                     │
│  UUID 已存在             IP 已存在                验证通过                    │
│  返回错误：               返回错误：              创建服务器记录               │
│  UUID已被注册             服务器已注册            状态为 active               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 详细设计

#### 1. Agent 端改造

**配置文件增加字段** `config/agent.yaml`：

```yaml
apiserver:
  http:
    scheme: http
    server: 127.0.0.1:10700
    baseUri: /api/v1
  # 新增配置
  registry:
    enabled: true           # 是否启用自动注册
    uuid: ""                # 预分配的 UUID（可选，为空则自动生成）
    token: ""               # 注册令牌（可选，用于安全验证）
```

**Agent 启动时注册逻辑** `internal/squ-agent/server/before_start.go`：

```go
func (s *Server) registerAgent() error {
    // 1. 检查是否已经注册过（读取本地存储的 server_id）
    serverID, err := s.getStoredServerID()
    if err == nil && serverID != "" {
        zap.L().Info("Agent already registered", 
            zap.String("server_id", serverID))
        return nil
    }
    
    // 2. 获取或生成 UUID
    uuid := s.Config.Apiserver.Registry.UUID
    if uuid == "" {
        uuid = generateUUID()  // agent 自己生成
    }
    
    // 3. 收集主机信息
    hostInfo, err := s.collectHostInfo()
    if err != nil {
        return fmt.Errorf("failed to collect host info: %w", err)
    }
    
    // 4. 向 apiserver 发起注册请求
    req := RegisterRequest{
        UUID:      uuid,
        Hostname:  hostInfo.Hostname,
        IPAddress: hostInfo.IPAddress,
        AgentPort: s.Config.Server.Port,
    }
    
    resp, err := s.callApiserverRegister(req)
    if err != nil {
        return fmt.Errorf("failed to register to apiserver: %w", err)
    }
    
    // 5. 注册成功，保存 server_id 到本地
    if err := s.storeServerID(resp.ServerID); err != nil {
        return fmt.Errorf("failed to store server_id: %w", err)
    }
    
    zap.L().Info("Agent registered successfully",
        zap.String("uuid", uuid),
        zap.Uint("server_id", resp.ServerID))
    
    return nil
}

// 检查是否已注册
func (s *Server) getStoredServerID() (string, error) {
    // 从 agent 本地数据库的 config 表读取
    conf, err := s.ConfRepository.Get("server_id")
    if err != nil {
        return "", err
    }
    return conf.Value, nil
}
```

#### 2. APIServer 端改造

**新增路由** `internal/squ-apiserver/router/server.go`：

```go
func Server(group *gin.RouterGroup, conf *config.Config, db database.DB) {
    // ... 现有路由 ...
    
    // 新增：Agent 注册路由（无需认证，或使用 token 验证）
    group.POST("/server/register", server.RegisterHandler(service))
}
```

**注册 Handler 改造** `internal/squ-apiserver/handler/server/service.go`：

```go
type RegisterRequest struct {
    UUID      string `json:"uuid" binding:"required"`
    Hostname  string `json:"hostname" binding:"required"`
    IPAddress string `json:"ip_address"`
    AgentPort int    `json:"agent_port" binding:"required"`
    Token     string `json:"token"`  // 可选的注册令牌
}

func (s *Server) Register(request req.Register) response.Response {
    // 1. 检查 UUID 是否已存在
    existingByUUID, err := s.Repository.GetByUUID(request.UUID)
    if err == nil {
        // UUID 已存在，检查是否是同一台服务器
        if existingByUUID.IpAddress == request.IPAddress {
            // 同一台服务器重新注册，更新端口即可
            existingByUUID.AgentPort = request.AgentPort
            s.Repository.Update(&existingByUUID)
            return response.Success(map[string]interface{}{
                "server_id": existingByUUID.ID,
                "status":    "updated",
            })
        }
        // 不同服务器尝试使用相同 UUID，拒绝
        return response.Error(res.ErrServerAlreadyRegistered)
    }
    
    // 2. 检查 IP 是否已存在
    existingByIP, err := s.Repository.GetByIPAddress(request.IPAddress)
    if err == nil {
        // IP 已存在但 UUID 不同，说明已经加入其他集群
        return response.Error(res.ErrServerAlreadyInCluster)
    }
    
    // 3. 创建新的服务器记录
    server := model.Server{
        UUID:      request.UUID,
        Hostname:  request.Hostname,
        IpAddress: request.IPAddress,
        AgentPort: request.AgentPort,
        Status:    model.ServerStatusOnline,
    }
    
    if err := s.Repository.Add(&server); err != nil {
        return response.Error(returnServerErrCode(err))
    }
    
    return response.Success(map[string]interface{}{
        "server_id": server.ID,
        "status":    "registered",
    })
}
```

#### 3. 错误码定义

`internal/squ-apiserver/handler/server/res/response_code.go`：

```go
const (
    // ... 现有错误码 ...
    ErrServerAlreadyRegistered  = 60007  // UUID 已被注册
    ErrServerAlreadyInCluster   = 60008  // 该服务器已加入其他集群
    ErrServerRegistrationFailed = 60009  // 注册失败
)
```

### 关键设计点

| 设计点 | 说明 |
|--------|------|
| **UUID 归属** | agent 生成或使用预分配 UUID，apiserver 只验证不生成 |
| **防重复加入** | 检查 UUID 和 IP 是否已存在，存在则拒绝 |
| **本地状态存储** | agent 本地存储 `server_id`，启动时先检查是否已注册 |
| **注册状态** | apiserver 的 Server 增加 `status` 字段：`pending`/`active`/`inactive` |
| **安全考虑** | 可选的 token 验证，防止恶意注册 |

### 两种加入方式的对比

| 方式 | 流程 | 适用场景 |
|------|------|----------|
| **从 APIServer 添加** | 先在页面创建 → 获取 UUID → 配置到 agent → agent 注册 | 需要预先规划的批量部署 |
| **Agent 主动加入** | 配置 APIServer 地址 → agent 启动时自动注册 | 快速部署、自动化场景 |

## 实现步骤

1. **Agent 端**
   - [ ] 在配置结构体中添加 `registry` 字段
   - [ ] 实现 `registerAgent()` 完整逻辑
   - [ ] 实现 `getStoredServerID()` 和 `storeServerID()` 方法
   - [ ] 实现 `collectHostInfo()` 方法收集主机信息
   - [ ] 实现 `callApiserverRegister()` 方法调用 APIServer 注册接口

2. **APIServer 端**
   - [ ] 在 Server model 中添加 `GetByUUID()` 和 `GetByIPAddress()` 方法
   - [ ] 完善 `Register()` 方法的实现
   - [ ] 添加注册路由 `/server/register`
   - [ ] 添加错误码定义
   - [ ] （可选）添加 token 验证机制

3. **前端**
   - [ ] 在服务器创建成功后显示 UUID，方便用户配置到 agent
   - [ ] 添加服务器状态显示（pending/active/inactive）

## 注意事项

1. **安全性**：注册接口无需认证，建议添加 token 验证机制防止恶意注册
2. **幂等性**：同一 agent 多次注册应该是幂等的，不会创建重复记录
3. **网络异常处理**：agent 注册失败时应该有重试机制
4. **配置热更新**：支持 agent 运行时修改 APIServer 地址
