# Agent 安全加入集群设计方案

## 背景

当前 Agent 与 APIServer 之间的通信存在以下问题：

1. **缺少身份验证**：Agent 与 APIServer 通信未使用双向认证
2. **证书分发困难**：需要手动在每台 Agent 机器上部署证书
3. **UUID 归属不明确**：UUID 生成位置不统一，Agent 无法确认自己的身份
4. **缺少防重复加入机制**：没有检查 Agent 是否已经加入其他集群

## 设计目标

1. **mTLS 双向认证（必须）**：Agent 与 APIServer 之间必须使用 mTLS 进行通信
2. **一次注册，永久使用**：Agent 注册成功后，自动获取证书，后续无需再次认证
3. **证书自动分发**：Agent 注册成功后，APIServer 自动签发并分发 Agent 专属证书
4. **防重复加入**：检查 Agent 是否已注册，避免重复注册

## 整体架构

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          Agent 安全加入流程                                    │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 生成 Join Token（Token 在有效期内可多次使用，适合 CI/CD 批量部署）          │
│  ┌─────────────┐                                                             │
│  │   Admin     │ ──► squctl token create ──► Token: abcdef.1234567890abcdef  │
│  └─────────────┘                                             │               │
│                                                              ▼               │
│  2. Agent 加入                                              复制 Token       │
│  ┌─────────────┐                                            到 Agent        │
│  │   Agent     │ ◄─────────────────────────────────────────────────          │
│  │  (工作节点)  │                                                             │
│  └──────┬──────┘                                                             │
│         │                                                                    │
│         │ POST /api/v1/agent/join                                            │
│         │ { token: "abcdef.1234567890abcdef", hostname, ip, ... }            │
│         │ [使用 Bootstrap 证书]                                               │
│         ▼                                                                    │
│  3. APIServer 处理                                                            │
│  ┌─────────────────┐                                                         │
│  │   APIServer     │ ──► 验证 Token 有效性（未过期）                           │
│  │                 │ ──► 检查 IP 是否已注册（防重复加入）                       │
│  │                 │ ──► 生成 Agent UUID + 专属证书                           │
│  │                 │ ──► 返回: UUID + CA证书 + Agent证书 + Agent私钥           │
│  └────────┬────────┘                                                         │
│           │                                                                  │
│           ▼                                                                  │
│  4. Agent 保存证书并启动                                                       │
│  ┌─────────────┐                                                             │
│  │   Agent     │ ──► 保存证书到本地                                           │
│  │             │ ──► 保存 UUID 到本地                                         │
│  │             │ ──► 后续通信使用 mTLS（无需 Token）                           │
│  └─────────────┘                                                             │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 详细设计

### 1. Join Token 设计

Token 在有效期内可**多次使用**，适合 CI/CD 批量部署场景。

```
格式: <token-id>.<token-secret>
示例: abcdef.0123456789abcdef

- token-id:     6 字符，用于 Token 查找
- token-secret: 16 字符，用于 Token 验证
```

#### Token 结构

```go
// internal/squ-apiserver/model/join_token.go

type JoinToken struct {
    ID        string    `json:"id"`          // Token ID (6字符)
    Secret    string    `json:"secret"`      // Token Secret (16字符，存储时加密)
    CreatedAt time.Time `json:"created_at"`  // 创建时间
    ExpiresAt time.Time `json:"expires_at"`  // 过期时间
    UsageCount int      `json:"usage_count"` // 使用次数（用于审计）
}

// Token 完整值（用于展示给用户）
func (t *JoinToken) FullToken() string {
    return fmt.Sprintf("%s.%s", t.ID, t.Secret)
}

// 是否有效
func (t *JoinToken) IsValid() bool {
    return time.Now().Before(t.ExpiresAt)
}
```

### 2. squctl 命令设计

#### Token 管理

```bash
# 创建默认 24 小时有效的 Token
squctl token create

# 创建指定有效期的 Token（适合 CI/CD 场景）
squctl token create --ttl 2h
squctl token create --ttl 168h  # 7 天，适合批量部署

# 输出示例
# Join token: abcdef.0123456789abcdef
# Expires: 2024-01-02 15:04:05 UTC

# 生成完整的加入命令
squctl token create --print-join-command

# 输出示例
# squirrel-agent join --apiserver https://apiserver.example.com:10700 --token abcdef.0123456789abcdef
```

#### Token 列表和删除

```bash
# 列出所有 Token
squctl token list

# 输出示例:
# ID       CREATED              EXPIRES              USAGE
# abcdef   2024-01-01 10:00:00  2024-01-02 10:00:00  5
# fedcba   2024-01-01 09:00:00  2024-01-01 10:00:00  12

# 删除 Token
squctl token delete <token-id>
```

### 3. APIServer 端实现

#### 新增路由

```go
// internal/squ-apiserver/router/agent.go

func Agent(group *gin.RouterGroup, conf *config.Config, db database.DB) {
    // Agent 加入接口（使用 Bootstrap 证书）
    group.POST("/agent/join", agent.JoinHandler(service))
    
    // 以下接口需要 mTLS 认证
    group.Use(mtls.Middleware(conf.MTLS))
    group.GET("/agent/certs/rotate", agent.RotateCertsHandler(service))
    group.GET("/agent/status", agent.StatusHandler(service))
}
```

#### Join Handler

```go
// internal/squ-apiserver/handler/agent/service.go

type JoinRequest struct {
    Token     string `json:"token" binding:"required"` // 完整 Token
    Hostname  string `json:"hostname" binding:"required"`
    IPAddress string `json:"ip_address"`
    AgentPort int    `json:"agent_port" binding:"required"`
}

type JoinResponse struct {
    ServerID    uint   `json:"server_id"`
    UUID        string `json:"uuid"`
    CACert      []byte `json:"ca_cert"`       // CA 证书 (PEM)
    AgentCert   []byte `json:"agent_cert"`    // Agent 证书 (PEM)
    AgentKey    []byte `json:"agent_key"`     // Agent 私钥 (PEM)
    APIServerCN string `json:"apiserver_cn"`  // APIServer 证书 CN（用于验证）
}

func (s *Service) Join(request JoinRequest) response.Response {
    // 1. 解析并验证 Token
    tokenID, tokenSecret := parseToken(request.Token)
    token, err := s.TokenRepo.GetByID(tokenID)
    if err != nil {
        return response.Error(ErrTokenNotFound)
    }
    
    // 2. 检查 Token 是否过期
    if !token.IsValid() {
        return response.Error(ErrTokenExpired)
    }
    
    // 3. 验证 Token Secret
    if !verifyTokenSecret(token.Secret, tokenSecret) {
        return response.Error(ErrTokenInvalid)
    }
    
    // 4. 检查 IP 是否已注册（防重复加入）
    if request.IPAddress != "" {
        existing, err := s.ServerRepo.GetByIPAddress(request.IPAddress)
        if err == nil && existing != nil {
            // IP 已存在，返回已有记录的 UUID，让 Agent 使用已有证书
            return response.Success(map[string]interface{}{
                "status":    "already_registered",
                "server_id": existing.ID,
                "uuid":      existing.UUID,
            })
        }
    }
    
    // 5. 生成 Agent UUID
    agentUUID := uuid.New().String()
    
    // 6. 生成 Agent 专属证书（CN = agent-{uuid}）
    certs, err := s.CertGenerator.GenerateAgentCert(agentUUID)
    if err != nil {
        return response.Error(ErrCertGenerationFailed)
    }
    
    // 7. 创建服务器记录
    server := &model.Server{
        UUID:      agentUUID,
        Hostname:  request.Hostname,
        IpAddress: request.IPAddress,
        AgentPort: request.AgentPort,
        Status:    model.ServerStatusOnline,
    }
    if err := s.ServerRepo.Add(server); err != nil {
        return response.Error(ErrServerCreateFailed)
    }
    
    // 8. 增加 Token 使用计数
    if err := s.TokenRepo.IncrementUsage(tokenID); err != nil {
        zap.L().Error("failed to increment token usage", zap.Error(err))
    }
    
    // 9. 返回证书和 UUID
    return response.Success(JoinResponse{
        ServerID:    server.ID,
        UUID:        agentUUID,
        CACert:      certs.CACert,
        AgentCert:   certs.AgentCert,
        AgentKey:    certs.AgentKey,
        APIServerCN: "squirrel-apiserver",
    })
}
```

#### 错误码定义

```go
// internal/squ-apiserver/handler/agent/res/response_code.go

const (
    ErrTokenNotFound         = 70001  // Token 不存在
    ErrTokenExpired          = 70002  // Token 已过期
    ErrTokenInvalid          = 70003  // Token 无效
    ErrServerAlreadyRegistered = 70004  // 服务器已注册
    ErrCertGenerationFailed  = 70005  // 证书生成失败
)
```

### 4. Agent 端实现

#### 配置文件

```yaml
# config/agent.yaml
apiserver:
  # 方式一：使用 Token 加入（首次启动）
  join:
    enabled: false
    server: "https://apiserver.example.com:10700"
    token: ""    # Join Token，格式: <id>.<secret>
  
  # 方式二：使用证书通信（已注册后）
  mtls:
    enabled: false
    caFile:   "/etc/squirrel/certs/ca.crt"
    certFile: "/etc/squirrel/certs/agent.crt"
    keyFile:  "/etc/squirrel/certs/agent.key"
    serverCN: "squirrel-apiserver"  # 服务端证书 CN，用于验证
```

#### Join 流程

```go
// internal/squ-agent/server/join.go

func (s *Server) joinCluster() error {
    // 1. 检查是否已有有效证书
    if s.hasValidCerts() {
        zap.L().Info("Agent already has valid certificates, skipping join")
        return nil
    }
    
    // 2. 检查是否已注册（本地存储的 UUID）
    if s.hasStoredUUID() {
        zap.L().Info("Agent already registered, checking certificates")
        // 尝试重新获取证书或使用本地缓存
        return s.verifyExistingRegistration()
    }
    
    // 3. 检查 Join 配置
    if !s.Config.Apiserver.Join.Enabled || s.Config.Apiserver.Join.Token == "" {
        return fmt.Errorf("join not configured, please provide join token")
    }
    
    // 4. 收集主机信息
    hostInfo, err := s.collectHostInfo()
    if err != nil {
        return fmt.Errorf("failed to collect host info: %w", err)
    }
    
    // 5. 发起 Join 请求
    req := JoinRequest{
        Token:     s.Config.Apiserver.Join.Token,
        Hostname:  hostInfo.Hostname,
        IPAddress: hostInfo.IPAddress,
        AgentPort: s.Config.Server.Port,
    }
    
    resp, err := s.callJoinAPI(req)
    if err != nil {
        return fmt.Errorf("join request failed: %w", err)
    }
    
    // 6. 保存证书
    certDir := "/etc/squirrel/certs"
    if err := s.saveCerts(certDir, resp); err != nil {
        return fmt.Errorf("failed to save certificates: %w", err)
    }
    
    // 7. 保存 UUID 到本地
    if err := s.storeUUID(resp.UUID); err != nil {
        return fmt.Errorf("failed to store UUID: %w", err)
    }
    
    // 8. 更新配置，启用 mTLS
    s.Config.Apiserver.MTLS.Enabled = true
    s.Config.Apiserver.MTLS.CAFile = filepath.Join(certDir, "ca.crt")
    s.Config.Apiserver.MTLS.CertFile = filepath.Join(certDir, "agent.crt")
    s.Config.Apiserver.MTLS.KeyFile = filepath.Join(certDir, "agent.key")
    
    zap.L().Info("Agent joined cluster successfully",
        zap.String("uuid", resp.UUID),
        zap.Uint("server_id", resp.ServerID))
    
    return nil
}

func (s *Server) saveCerts(dir string, resp JoinResponse) error {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    // 保存 CA 证书
    if err := os.WriteFile(filepath.Join(dir, "ca.crt"), resp.CACert, 0644); err != nil {
        return err
    }
    
    // 保存 Agent 证书
    if err := os.WriteFile(filepath.Join(dir, "agent.crt"), resp.AgentCert, 0644); err != nil {
        return err
    }
    
    // 保存 Agent 私钥（权限更严格）
    if err := os.WriteFile(filepath.Join(dir, "agent.key"), resp.AgentKey, 0600); err != nil {
        return err
    }
    
    return nil
}

// 检查本地是否已有 UUID
func (s *Server) hasStoredUUID() bool {
    uuid, err := s.getStoredUUID()
    return err == nil && uuid != ""
}

func (s *Server) getStoredUUID() (string, error) {
    conf, err := s.ConfRepository.Get("agent_uuid")
    if err != nil {
        return "", err
    }
    return conf.Value, nil
}

func (s *Server) storeUUID(uuid string) error {
    return s.ConfRepository.Set("agent_uuid", uuid)
}
```

### 5. 证书生成器

```go
// internal/squctl/certs/agent.go

// GenerateAgentCert 为 Agent 生成专属证书
// CN 格式: agent-{uuid}
func (g *Generator) GenerateAgentCert(uuid string) (*AgentCerts, error) {
    // 确保 CA 可用
    if err := g.ensureCA(); err != nil {
        return nil, err
    }
    
    // 生成私钥
    key, err := rsa.GenerateKey(rand.Reader, g.keySize)
    if err != nil {
        return nil, fmt.Errorf("failed to generate agent key: %w", err)
    }
    
    // 创建证书模板
    cn := fmt.Sprintf("agent-%s", uuid)
    serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    
    template := &x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            CommonName:         cn,
            Organization:       []string{"Squirrel"},
            OrganizationalUnit: []string{"Agent"},
        },
        NotBefore:   time.Now(),
        NotAfter:    time.Now().Add(g.expiry),
        KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
    }
    
    // 使用 CA 签名
    certDER, err := x509.CreateCertificate(rand.Reader, template, g.caCert, &key.PublicKey, g.caKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create agent certificate: %w", err)
    }
    
    // 读取 CA 证书
    caCertPEM, _ := os.ReadFile(filepath.Join(g.outputDir, "ca.crt"))
    
    return &AgentCerts{
        CACert:    caCertPEM,
        AgentCert: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}),
        AgentKey:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}),
    }, nil
}

type AgentCerts struct {
    CACert    []byte
    AgentCert []byte
    AgentKey  []byte
}
```

### 6. mTLS 中间件

```go
// internal/squ-apiserver/middleware/mtls.go

func Middleware(cfg *config.MTLS) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 验证客户端证书
        if c.Request.TLS == nil || len(c.Request.TLS.PeerCertificates) == 0 {
            c.AbortWithStatusJSON(403, gin.H{"error": "client certificate required"})
            return
        }
        
        cert := c.Request.TLS.PeerCertificates[0]
        
        // 2. 验证证书 CN
        cn := cert.Subject.CommonName
        if !isValidCN(cn, cfg.AllowedCNs) {
            c.AbortWithStatusJSON(403, gin.H{"error": "unauthorized client certificate"})
            return
        }
        
        // 3. 将 Agent UUID 存入上下文
        if strings.HasPrefix(cn, "agent-") {
            agentUUID := strings.TrimPrefix(cn, "agent-")
            c.Set("agent_uuid", agentUUID)
            c.Set("is_agent", true)
        }
        
        c.Next()
    }
}

func isValidCN(cn string, allowedCNs []string) bool {
    for _, allowed := range allowedCNs {
        if cn == allowed {
            return true
        }
    }
    // Agent 证书以 "agent-" 开头
    return strings.HasPrefix(cn, "agent-")
}
```

## 使用流程

### 管理员操作

```bash
# 1. 启动 APIServer（证书配置在配置文件中）
squirrel-apiserver

# 2. 生成 Join Token（可设置较长有效期用于 CI/CD 批量部署）
squctl token create --ttl 168h --print-join-command

# 输出:
# squirrel-agent join --apiserver https://192.168.1.10:10700 --token abcdef.0123456789abcdef
```

### Agent 加入集群

```bash
# 方式一：命令行加入
squirrel-agent join --apiserver https://192.168.1.10:10700 --token abcdef.0123456789abcdef

# 方式二：配置文件方式
# 编辑 /etc/squirrel/agent.yaml
# 
# apiserver:
#   join:
#     enabled: true
#     server: "https://192.168.1.10:10700"
#     token: "abcdef.0123456789abcdef"

# 启动 Agent
squirrel-agent
```

### CI/CD 批量部署示例

```bash
# 创建一个长有效期 Token（7天）
TOKEN=$(squctl token create --ttl 168h | grep "Join token:" | awk '{print $3}')

# 在 CI/CD 中使用该 Token 批量部署 Agent
for host in $(cat hosts.txt); do
    ssh $host "squirrel-agent join --apiserver https://apiserver.example.com:10700 --token $TOKEN"
done
```

## 安全设计

| 安全点 | 说明 |
|--------|------|
| **mTLS 强制** | 所有 Agent 与 APIServer 通信必须使用 mTLS |
| **Token 有效期** | 默认 24 小时，最长建议不超过 7 天 |
| **Token 可复用** | 有效期内可多次使用，适合 CI/CD 批量部署 |
| **证书归属** | 每个 Agent 获得专属证书，CN 包含 UUID |
| **私钥保护** | Agent 私钥权限设为 0600，仅 root 可读 |
| **防重复加入** | 检查 IP 是否已注册，已注册则返回已有 UUID |

## Agent 状态管理

### 本地状态存储

Agent 在本地数据库存储以下信息：

| 字段 | 说明 |
|------|------|
| `agent_uuid` | Agent 的唯一标识，由 APIServer 分配 |
| `server_id` | 在 APIServer 中的记录 ID |
| `join_time` | 加入集群的时间 |

### 状态检查流程

```
Agent 启动
    │
    ├─── 检查本地是否有 UUID
    │        │
    │        ├─── 有 ──► 检查证书是否有效
    │        │              │
    │        │              ├─── 有效 ──► 启动 mTLS 通信
    │        │              │
    │        │              └─── 无效 ──► 尝试证书轮换
    │        │
    │        └─── 无 ──► 检查是否有 Join Token
    │                         │
    │                         ├─── 有 ──► 发起 Join 请求
    │                         │
    │                         └─── 无 ──► 报错退出
```

## 异常处理

### Agent 离线重连

- Agent 重启时，检查本地 UUID 和证书
- 证书有效则直接使用 mTLS 连接
- 证书即将过期，主动请求证书轮换

### Agent 被删除后重新加入

- APIServer 删除服务器记录后，该 Agent 的证书失效
- Agent 需要重新使用 Token 加入

### 证书过期处理

- Agent 主动检测证书有效期
- 证书即将过期（< 7 天），自动请求轮换
- 轮换失败则告警

## 实现步骤

### Phase 1: Token 机制

- [ ] 实现 `JoinToken` model 和 repository
- [ ] 实现 `squctl token create/list/delete` 命令
- [ ] 实现 APIServer `/agent/join` 接口
- [ ] 实现 Agent Join 流程

### Phase 2: 证书管理

- [ ] 扩展 `GenerateAgentCert(uuid)` 方法
- [ ] Agent 保存证书到本地
- [ ] mTLS 中间件实现
- [ ] Agent 本地 UUID 存储

### Phase 3: 运维增强

- [ ] 实现证书轮换 API
- [ ] 实现证书过期告警
- [ ] Token 使用统计和审计
- [ ] 完善日志和监控

## 参考设计

- [kubeadm join 工作原理](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/)
- [Bootstrap Tokens](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/)
